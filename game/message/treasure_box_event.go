package message

import (
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/proto"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

//宝箱事件分发
type TreasureBoxEvent struct {
}

var (
	G_TreasureBoxRandMap               sync.Map    //随机地图，key1城市编号,key2坐标编号 value 事件数据
	G_TreasureBoxEventConf             sync.Map    //key 事件序号，value事件信息
	G_TreasureBoxEventRunning          bool        = false
	G_TreasureBoxEventRunningTimestamp int64       = 0                                                                   //宝箱事件启动的时间戳
	G_TreasureBoxTimerReset            int         = 9999                                                                //定时器时间间隔
	G_TreasureBoxTimer                 *time.Timer = time.NewTimer(time.Duration(G_TreasureBoxTimerReset) * time.Second) //定时器
	G_TreasureBox_Clean_Map_Flag       bool        = false
	G_TreasureBox_Operation_Lock       sync.Mutex  // 宝箱操作锁
)

const (
	TimerTreasureBoxName = "TreasureBoxEvent"
	TreasureBox_MAP      = "treasure_box_map"
)

// 设置配置信息
func TreasureBoxEvent_Config(cfg *proto.BaseConf) {
	for _, v := range cfg.TreasureBoxEventRate {
		G_TreasureBoxEventConf.Store(v.Id, v)
	}
}

// 初始化地图
func (this *TreasureBoxEvent) InitMap(locationId int, clean bool) {
	logger.Debugf("TreasureBoxEvent_SetMap in")

	value, ok := G_TreasureBoxRandMap.Load(locationId)
	if !clean && ok {
		node := value.(*sync.Map)

		allValue, ok := G_CityAllMap.Load(locationId)
		if !ok {
			var emptyMap sync.Map
			G_TreasureBoxRandMap.Store(locationId, &emptyMap)

			logger.Debugf("TreasureBoxEvent_SetMap end")
			return
		}
		allMap := allValue.(*sync.Map)
		node.Range(func(key, value interface{}) bool {
			// logger.Errorf("TreasureBoxEvent_InitMap values:", key, value)
			allMap.Delete(key)
			return true
		})
	} else {
		// 清理时发通知给前端清理事件
		if clean == true {
			if ok {
				node := value.(*sync.Map)
				events := this.GetMapEvents(node)
				if events != nil && len(events) > 0 {
					response := &proto.S2CBroadRandEvent{}
					response.LocationID = locationId
					response.EventInfo = events

					packet := &utils.Packet{}
					packet.Initialize(proto.MSG_BROAD_CLEAN_EVENT)
					packet.WriteData(response)

					logger.Debugf("TreasureBoxEvent BROAD_CLEAN_EVENT %+V %+V", locationId, events)
					Sched.BroadCastMsgToPlatform(int32(locationId), "", packet, proto.PLATFORM_ANDROID)
				}
			}
		}
		var emptyMap sync.Map
		G_TreasureBoxRandMap.Store(locationId, &emptyMap)
	}
	logger.Debugf("TreasureBoxEvent_SetMap end")
}

// 获取事件
func TreasureBoxEvent_GetEvent(locationId int, positionNo int) *proto.EventData {
	if value, ok := G_TreasureBoxRandMap.Load(locationId); ok {
		node := value.(*sync.Map)
		if subvalue, ok := node.Load(positionNo); ok {
			ret := subvalue.(*proto.EventData)
			return ret
		}
	}

	return nil
}

// 重置事件
func TreasureBoxEvent_ResetEvent(locationId int, positionNo int, eventData *proto.EventData) {
	// logger.Debugf("TreasureBoxEvent_SetMap in")

	if value, ok := G_TreasureBoxRandMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Delete(positionNo)
		// G_TreasureBoxRandMap.Store(locationId, node)
	}

	if value, ok := G_CityAllMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Store(positionNo, eventData.EventNode)
	}
	// logger.Debugf("TreasureBoxEvent_SetMap end")
}

// 查找地图节点
func (this *TreasureBoxEvent) GetMap(locationId int) *sync.Map {
	value, ok := G_CityAllMap.Load(locationId)
	if !ok {
		return nil
	}

	nodeAll := value.(*sync.Map)
	return nodeAll
}

func (this *TreasureBoxEvent) Start() {
	G_TreasureBox_Operation_Lock.Lock()
	defer G_TreasureBox_Operation_Lock.Unlock()

	status := G_ActivityManage.GetActivityStatus(proto.ACTIVITY_TREASURE_BOX)
	if status != proto.ACTIVITY_START {
		G_BaseCfg.TreasureBox.StateSwitch = 0
		if G_TreasureBoxEventRunning {
			G_TreasureBoxTimer.Stop()
			G_TreasureBoxEventRunning = false
		}
	} else {
		G_TreasureBox_Clean_Map_Flag = false
		G_BaseCfg.TreasureBox.StateSwitch = 1
		G_TreasureBoxEventRunning = true
		G_TreasureBoxEventRunningTimestamp = time.Now().Unix()
		G_TreasureBoxTimerReset = G_BaseCfg.TreasureBox.EventInterval
		G_TreasureBoxTimer.Reset(time.Duration(1) * time.Second)
	}
	//广播消息活动状态
	go PushActivityStatus(true, "")
}

func (this *TreasureBoxEvent) Stop() {
	G_TreasureBox_Operation_Lock.Lock()
	defer G_TreasureBox_Operation_Lock.Unlock()

	last_TreasureBoxEventRunning := G_TreasureBoxEventRunning
	G_TreasureBoxTimer.Stop()

	G_TreasureBoxEventRunning = false
	if last_TreasureBoxEventRunning {
		//广播消息
		go PushActivityStatus(true, "")
	}
}

func (this *TreasureBoxEvent) StopAndClean() {
	G_TreasureBox_Operation_Lock.Lock()
	defer G_TreasureBox_Operation_Lock.Unlock()

	last_TreasureBoxEventRunning := G_TreasureBoxEventRunning
	G_TreasureBoxTimer.Stop()

	G_TreasureBoxEventRunning = false
	G_TreasureBox_Clean_Map_Flag = true
	this.CleanRandMap()

	if last_TreasureBoxEventRunning {
		//广播消息
		go PushActivityStatus(true, "")
	}
}

func (this *TreasureBoxEvent) GetMapEvents(nodeMap *sync.Map) map[int]*proto.EventData {
	if nodeMap == nil {
		return nil
	}

	events := make(map[int]*proto.EventData)
	nodeMap.Range(func(key interface{}, value interface{}) bool {
		positionId := key.(int)
		event_data := value.(*proto.EventData)

		events[positionId] = event_data
		return true
	})
	return events
}

func (this *TreasureBoxEvent) CleanRandMap() {
	for _, locationId := range G_BaseCfg.LocationId {
		value, ok := G_TreasureBoxRandMap.Load(locationId)
		if !ok {
			logger.Errorf("G_TreasureBoxRandMap.Load(locationId) failed", locationId)
			continue
		}

		nodeRand := value.(*sync.Map)
		// 待删除的事件
		removedEventData := this.RemoveRandMap(locationId, nodeRand)

		// 广播清理掉过期事件给前端
		if removedEventData != nil && len(removedEventData) > 0 {
			response := &proto.S2CBroadRandEvent{}
			response.LocationID = locationId
			response.EventInfo = removedEventData

			packet := &utils.Packet{}
			packet.Initialize(proto.MSG_BROAD_CLEAN_EVENT)
			packet.WriteData(response)

			logger.Debugf("TreasureBoxEvent BROAD_CLEAN_EVENT %+V %+V", locationId, removedEventData)
			Sched.BroadCastMsgToPlatform(int32(locationId), "", packet, proto.PLATFORM_ANDROID)
		}
	}
}

func (this *TreasureBoxEvent) RemoveRandMap(locationId int, nodeRand *sync.Map) map[int]*proto.EventData {
	if nodeRand == nil {
		return nil
	}

	removed_event := make(map[int]*proto.EventData)
	nodeRand.Range(func(key interface{}, value interface{}) bool {
		positionId := key.(int)
		event_data := value.(*proto.EventData)

		TreasureBoxEvent_ResetEvent(locationId, positionId, event_data)
		removed_event[positionId] = event_data
		return true
	})
	return removed_event
}

func (this *TreasureBoxEvent) EventProcess() {
	status := G_ActivityManage.GetActivityStatus(proto.ACTIVITY_TREASURE_BOX)
	// 关闭活动
	if status == 0 {
		G_BaseCfg.TreasureBox.StateSwitch = 0
		go PushActivityStatus(true, "")
		return
	}
	G_TreasureBoxTimer.Reset(time.Duration(G_TreasureBoxTimerReset) * time.Second)
	// logger.Debugf("TreasureBoxEvent EventProcess ", time.Now(), G_TreasureBoxEventRunning)
	if !G_TreasureBoxEventRunning {
		return
	}
	this.BrocastEvent()
}

//产生并广播随机事件
func (this *TreasureBoxEvent) BrocastEvent() {

	for _, locationId := range G_BaseCfg.LocationId {
		value, ok := G_TreasureBoxRandMap.Load(locationId)
		if !ok {
			logger.Errorf("G_TreasureBoxRandMap.Load(locationId) failed", locationId)
			continue
		}

		nodeRand := value.(*sync.Map)
		// 待删除的事件
		removedEventData := this.checkExpiredAndRemove(locationId, nodeRand)

		// 广播清理掉过期事件给前端
		if removedEventData != nil && len(removedEventData) > 0 {
			response := &proto.S2CBroadRandEvent{}
			response.LocationID = locationId
			response.EventInfo = removedEventData

			packet := &utils.Packet{}
			packet.Initialize(proto.MSG_BROAD_CLEAN_EVENT)
			packet.WriteData(response)

			logger.Debugf("TreasureBoxEvent BROAD_CLEAN_EVENT %+V %+V", locationId, removedEventData)
			Sched.BroadCastMsgToPlatform(int32(locationId), "", packet, proto.PLATFORM_ANDROID)
		}

		eventData := this.getRandPositionMap(locationId, nodeRand)

		// 分发新增事件给前端
		if len(eventData) == 0 {
			continue
		}
		for k, v := range eventData {
			// if _, ok := nodeRand.Load(k); ok {
			// 	logger.Error("nodeRand exists key", k)
			// }
			nodeRand.Store(k, v)
		}
		// 清理所有的事件
		if G_TreasureBox_Clean_Map_Flag {
			this.RemoveRandMap(locationId, nodeRand)
		}
		// dataMap := this.getRandMap()
		// nodeRand, ok := dataMap[locationId]
		// if !ok {
		// 	continue
		// }
		// buf, _ := json.Marshal(nodeRand)
		// utils.SetKeyValue(TreasureBox_MAP, strconv.Itoa(locationId), buf, false, utils.ITEM_DAY)

		response := &proto.S2CBroadRandEvent{}
		response.LocationID = locationId
		response.EventInfo = eventData

		packet := &utils.Packet{}
		packet.Initialize(proto.MSG_BROAD_RAND_EVENT)
		packet.WriteData(response)

		// logger.Debugf("TreasureBoxEvent brocastEvent %+V %+V", locationId, eventData)

		Sched.BroadCastMsgToPlatform(int32(locationId), "", packet, proto.PLATFORM_ANDROID)
	}
}

// 获取随机地图元素
func (this *TreasureBoxEvent) getRandMap() map[int]map[int]*proto.EventData {
	ret := make(map[int]map[int]*proto.EventData)
	G_TreasureBoxRandMap.Range(func(k interface{}, v interface{}) bool {
		locationId := k.(int)
		ret[locationId] = make(map[int]*proto.EventData, 0)
		data := v.(*sync.Map)
		data.Range(func(key interface{}, value interface{}) bool {
			positionId := key.(int)
			ret[locationId][positionId] = value.(*proto.EventData)
			return true
		})
		return true
	})
	return ret
}

// 获取地图元素
func (this *TreasureBoxEvent) getMapKeys(node *sync.Map) map[int]int {
	ret := make(map[int]int, 0)
	i := 1
	node.Range(func(key interface{}, value interface{}) bool {
		ret[i] = key.(int)
		i++
		return true
	})
	return ret
}

//根据指定的源map随机取指定数量的元素,并返回取的源map的元素集合
func (this *TreasureBoxEvent) getRandPositionMap(locationId int, nodeRand *sync.Map) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData, 0)

	if len(G_BaseCfg.TreasureBoxEventRate) == 0 {
		logger.Errorf("G_BaseCfg.TreasureBoxEventRate failed", G_BaseCfg.TreasureBoxEventRate)
		return ret
	}

	mapAll := this.GetMap(locationId)
	if mapAll == nil {
		return ret
	}

	G_TreasureBoxEventConf.Range(func(key interface{}, value interface{}) bool {
		nodeEvent := value.(proto.EventRate)
		eventData := this.generateEvent(locationId, mapAll, nodeRand, &nodeEvent)
		if len(eventData) > 0 {
			for k, v := range eventData {
				ret[k] = v
			}
		}
		return true
	})

	return ret
}

// 检查组件是否锁定
func (this *TreasureBoxEvent) checkExpiredEventIsLocked(positionNo int) bool {
	key := BOX_IS_LOCK + strconv.Itoa(positionNo)
	res, err := db.RedisMgr.GetRedisClient().Exists(key).Result()
	// logger.Debugf("TreasureBoxEvent checkExpiredEventIsLocked %+V %+V", res, err)
	if err != nil {
		return false
	}
	if res == 0 {
		return false
	}
	return true
}

// 检查过期并删除
func (this *TreasureBoxEvent) checkExpiredAndRemove(locationId int, nodeRand *sync.Map) map[int]*proto.EventData {
	if G_BaseCfg.TreasureBox.EventTimeoutSwitch == 0 {
		return nil
	}

	eventLen := utils.GetMapLen(nodeRand)
	if eventLen <= 0 {
		return nil
	}
	now_timestamp := time.Now().Unix()

	if int(now_timestamp-G_TreasureBoxEventRunningTimestamp) < G_BaseCfg.TreasureBox.EventTimeoutCheckInterval {
		return nil
	} else {
		G_TreasureBoxEventRunningTimestamp = now_timestamp
	}

	removed_event := make(map[int]*proto.EventData)
	nodeRand.Range(func(key interface{}, value interface{}) bool {
		positionId := key.(int)
		event_data := value.(*proto.EventData)
		expired_event := (*proto.ExpiredEventData)(unsafe.Pointer(event_data))
		if int(now_timestamp-expired_event.Timestamp) >= G_BaseCfg.TreasureBox.EventTimeout {
			if !this.checkExpiredEventIsLocked(positionId) {
				removed_event[positionId] = event_data
			}
		}
		return true
	})
	for positionId, event_data := range removed_event {
		TreasureBoxEvent_ResetEvent(locationId, positionId, event_data)
	}
	return removed_event
}

// 生成事件
func (this *TreasureBoxEvent) generateEvent(locationId int, nodeAll *sync.Map, nodeRand *sync.Map, eventRate *proto.EventRate) map[int]*proto.EventData {
	response := make(map[int]*proto.EventData)

	eventLen := utils.GetMapLen(nodeRand)
	generateNum := G_BaseCfg.TreasureBox.EventNum - eventLen
	// logger.Debugf("generateEvent ", eventLen, G_BaseCfg.TreasureBox.EventNum, generateNum)
	if generateNum <= 0 {
		return response
	}

	for i := 0; i < generateNum; i++ {
		leftMap := this.getMapKeys(nodeAll)
		lenLeft := len(leftMap)

		if lenLeft == 0 {
			break
		}

		idx := rand.Intn(lenLeft) + 1
		positionNo := leftMap[idx]
		var node interface{}
		var ok bool
		if node, ok = nodeAll.Load(positionNo); !ok {
			break
		}
		eventNode := node.(proto.EventNode)

		event_data := new(proto.ExpiredEventData)
		event_data.ActivityType = proto.ACTIVITY_TYPE_TREASURE_BOX
		event_data.Type = eventNode.Type
		event_data.X = eventNode.X
		event_data.Y = eventNode.Y
		event_data.EventId = eventRate.Id
		event_data.Timestamp = time.Now().Unix()

		out_evnet_data := (*proto.EventData)(unsafe.Pointer(event_data))
		// logger.Errorf("generateEvent positionNo,", positionNo)
		nodeAll.Delete(positionNo)

		response[positionNo] = out_evnet_data
	}
	return response
}
