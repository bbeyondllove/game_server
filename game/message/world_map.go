package message

import (
	"encoding/json"
	"game_server/core/utils"
	"game_server/game/message/activity/red_envelope"

	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/proto"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"game_server/core/logger"

	"game_server/core/base"

	"github.com/shopspring/decimal"
)

//普通事件分发
type NormalEvent struct {
}

var (
	G_BaseCfg               = new(proto.BaseConf) //配置信息
	G_CityAllMap            sync.Map              //所有地图，key1城市编号,key2坐标编号 value "x,y"坐标
	G_CityRandMap           sync.Map              //随机地图，key1城市编号,key2坐标编号 value 事件数据
	G_EventConf             sync.Map              //key 事件序号，value事件信息
	G_UndoTaskMap           sync.Map              //未完成任务列表
	G_DoubleYearEvent                             = new(DoubleYearEvent)
	G_TreasureBoxEvent                            = new(TreasureBoxEvent)
	G_NormalEventTimerReset int                   = 9999                                                                //定时器时间间隔
	G_NormalEventTimer      *time.Timer           = time.NewTimer(time.Duration(G_NormalEventTimerReset) * time.Second) //定时器
	G_NormalEventRunning    bool                  = false
	G_NormalEvent                                 = new(NormalEvent)
)

func InitDbData() {
	decimal.DivisionPrecision = 2
	NickNameInit()
	setUserTask()
	ReloadDb()

	go func() {
		for {
			err := reflushStoreInfo()
			if err != nil {
				logger.Errorf("refulsh reflushStoreInfo failed:", err.Error())
			}
			time.Sleep(_SleepTime)

		}
	}()
}

func ReloadDb() {
	logger.Debugf("ReloadDb in")
	go WorldMapInit(false)
	logger.Debugf("ReloadDb in")
}

func reflushStoreInfo() error {
	GetItemList()

	buildingTypeInfos, err := db_service.BuildTypeIns.GetBuildingTypes()
	if err != nil {
		logger.Errorf("reflush GetBuildingTypes() failed:", err.Error())
		return err
	}
	BuildingTypes.Range(func(key interface{}, value interface{}) bool {
		BuildingTypes.Delete(key)
		return true
	})
	for _, info := range buildingTypeInfos {
		BuildingTypes.Store(info.SmallType, info)
	}

	return nil
}

//城市事件初始化
func WorldMapInit(clean_rand_map bool) {
	logger.Debugf("WorldMapInit in")
	G_CityAllMap.Range(func(key interface{}, value interface{}) bool {
		G_CityAllMap.Delete(key)
		return true
	})

	//双蛋地图初始化
	G_DoubleYearEvent.EventInit()
	G_NormalEvent.Stop()
	G_TreasureBoxEvent.Stop()

	config := G_ActivityManage.GetActivityConfig(proto.ACTIVITY_SPRING_FESTIVAL)
	if config != nil {
		base.Setting.Springfestival.ActivityStartDatetime = utils.Time2Str(config.StartTime)
		base.Setting.Springfestival.ActivityEndDatetime = utils.Time2Str(config.FinishTime)
	}

	status := G_DoubleYearEvent.getActivityState()

	//读配置文件,并启动定时器生成事件
	for _, v := range G_BaseCfg.LocationId {

		//mapAry, err := db_service.WorldMapIns.GetAllBuilding(v)
		//if err != nil || len(mapAry) == 0 {
		//	continue
		//}

		var dataNormal sync.Map
		var dataDouble sync.Map

		for k, value := range G_BaseCfg.EventMap {
			positionNo := getPositionNo(value.X, value.Y)
			var cityEvent proto.EventNode
			//cityEvent.Type = value.SmallType
			cityEvent.X = value.X
			cityEvent.Y = value.Y
			cityEvent.ActivityType = proto.ACTIVITY_TYPE_NOMAL

			////加上双旦判断
			//if status != proto.ACTIVITY_END &&
			//	cityEvent.X >= G_BaseCfg.DoubleYearArea["start_position"].X &&
			//	cityEvent.X <= G_BaseCfg.DoubleYearArea["end_position"].X &&
			//	cityEvent.Y >= G_BaseCfg.DoubleYearArea["start_position"].Y &&
			//	cityEvent.Y <= G_BaseCfg.DoubleYearArea["end_position"].Y {
			//	cityEvent.ActivityType = proto.ACTIVITY_TYPE_DOUBLE_YEAR
			//	dataDouble.Store(positionNo, cityEvent)
			//
			//} else {
			//	dataNormal.Store(positionNo, cityEvent)
			//}

			if status != proto.ACTIVITY_END {
				if k%5 <= 2 {
					cityEvent.ActivityType = proto.ACTIVITY_TYPE_DOUBLE_YEAR
					dataDouble.Store(positionNo, cityEvent)
				} else {
					dataNormal.Store(positionNo, cityEvent)
				}
			} else {
				dataNormal.Store(positionNo, cityEvent)
			}
		}
		G_CityAllMap.Store(v, &dataNormal)
		Sched.InitConnMap(int32(v))

		G_DoubleYearEvent.EventStore(v, dataDouble)

		G_NormalEvent.InitMap(v, clean_rand_map)
		G_TreasureBoxEvent.InitMap(v, clean_rand_map)
	}

	G_NormalEvent.Start()
	G_TreasureBoxEvent.Start()

	G_DoubleYearEvent.loadMap()
	G_DoubleYearEvent.Start()

	logger.Debugf("WorldMapInit end")
}

//读取base.json基础配置文件
func ReadCfg(cfg *proto.BaseConf) bool {
	//ReadFile函数会读取文件的全部内容，并将结果以[]
	filename := GetExecpath() + "/base.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("读取配置文件错误")
		return false
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		logger.Errorf(" parse json error ", err)
		return false
	}

	for _, v := range cfg.EventRate {
		G_EventConf.Store(v.Id, v)
	}

	for _, v := range cfg.DoubleEventRate {
		G_ChristmaEventConf.Store(v.Id, v)
	}

	TreasureBoxEvent_Config(cfg)
	TreasureBox_Config(cfg)
	red_envelope.RedEnvelopeConfig()

	return true
}

// 获取当前程序运行目录
func GetExecpath() string {
	execpath, _ := os.Executable() // 获得程序路径
	path := filepath.Dir(execpath)
	return strings.Replace(path, "\\", "/", -1)
}

//产生并广播随机事件
func (this *NormalEvent) brocastEvent() {
	// logger.Debugf("brocastEvent in")

	for _, v := range G_BaseCfg.LocationId {
		value, ok := G_CityRandMap.Load(v)
		if !ok {
			continue
		}
		nodeRand := value.(*sync.Map)

		eventData := this.getRandPositionMap(v, nodeRand)
		if len(eventData) == 0 {
			continue
		}
		for k, v := range eventData {
			if _, ok := nodeRand.Load(k); ok {
				logger.Errorf("brocastEvent exists key", k)
			}
			nodeRand.Store(k, v)
		}
		rsp := &utils.Packet{}
		rsp.Initialize(proto.MSG_BROAD_RAND_EVENT)
		responseMessage := &proto.S2CBroadRandEvent{}
		responseMessage.LocationID = v
		responseMessage.EventInfo = eventData
		rsp.WriteData(responseMessage)
		Sched.BroadCastMsg(int32(v), "", rsp)
		logger.Debugf("brocastEvent retEvent=%+V %+V", v, eventData)

	}
	// logger.Debugf("brocastEvent end")
}

//获取第index个事件id
func getEventId(index int) int {
	i := 0
	ret := 0
	G_EventConf.Range(func(key interface{}, value interface{}) bool {
		if i == index {
			ret = key.(int)
			return true
		} else if i < index {
			i++
		}
		return true
	})

	return ret
}

func getPositionNo(x, y int) int {
	return x + G_BaseCfg.EventOffset*y
}

func getEventData(activityType, locationId, x, y int) *proto.EventData {
	positionNo := getPositionNo(x, y)
	if activityType == proto.ACTIVITY_TYPE_NOMAL {
		if value, ok := G_CityRandMap.Load(locationId); ok {
			node := value.(*sync.Map)
			if subvalue, ok := node.Load(positionNo); ok {
				ret := subvalue.(*proto.EventData)
				return ret
			}
		}
	} else if activityType == proto.ACTIVITY_TYPE_DOUBLE_YEAR {
		if value, ok := G_ChristmaRandMap.Load(locationId); ok {
			node := value.(*sync.Map)
			if subvalue, ok := node.Load(positionNo); ok {
				ret := subvalue.(*proto.EventData)
				return ret
			}
		}
	} else if activityType == proto.ACTIVITY_TYPE_TREASURE_BOX {
		return TreasureBoxEvent_GetEvent(locationId, positionNo)
	}

	return nil
}

func getAreaEventData(locationId int, allow_treasureBox bool) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData)
	if value, ok := G_CityRandMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Range(func(key interface{}, value interface{}) bool {
			ret[key.(int)] = value.(*proto.EventData)
			return true
		})
	}

	if allow_treasureBox {
		if value, ok := G_TreasureBoxRandMap.Load(locationId); ok {
			node := value.(*sync.Map)
			node.Range(func(key interface{}, value interface{}) bool {
				ret[key.(int)] = value.(*proto.EventData)
				return true
			})
		}
	}

	if value, ok := G_ChristmaRandMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Range(func(key interface{}, value interface{}) bool {
			ret[key.(int)] = value.(*proto.EventData)
			return true
		})
	}

	return ret
}

func resetEventData(activityType, locationId int, positionNo int, eventData *proto.EventData) {
	if activityType == proto.ACTIVITY_TYPE_NOMAL {
		if value, ok := G_CityAllMap.Load(locationId); ok {
			node := value.(*sync.Map)
			node.Store(positionNo, eventData.EventNode)
			G_CityAllMap.Store(locationId, node)
		}
		if value, ok := G_CityRandMap.Load(locationId); ok {
			node := value.(*sync.Map)
			node.Delete(positionNo)
			G_CityRandMap.Store(locationId, node)
		}
	} else if activityType == proto.ACTIVITY_TYPE_DOUBLE_YEAR {
		if value, ok := G_ChristmaAllMap.Load(locationId); ok {
			node := value.(*sync.Map)
			node.Store(positionNo, eventData.EventNode)
			G_ChristmaAllMap.Store(locationId, node)
		}
		if value, ok := G_ChristmaRandMap.Load(locationId); ok {
			node := value.(*sync.Map)
			node.Delete(positionNo)
			G_ChristmaRandMap.Store(locationId, node)
		}
	} else if activityType == proto.ACTIVITY_TYPE_TREASURE_BOX {
		TreasureBoxEvent_ResetEvent(locationId, positionNo, eventData)
	}

}

func getLeftMap(locationId int) map[int]int {
	ret := make(map[int]int, 0)
	i := 1
	if value, ok := G_CityAllMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Range(func(key interface{}, value interface{}) bool {
			ret[i] = key.(int)
			i++
			return true
		})
	}

	return ret
}

func (this *NormalEvent) InitMap(locationId int, clean bool) {
	logger.Debugf("NormalEvent in")

	value, ok := G_CityRandMap.Load(locationId)
	if !clean && ok {
		node := value.(*sync.Map)

		allValue, ok := G_CityAllMap.Load(locationId)
		if !ok {
			var emptyMap sync.Map
			G_CityRandMap.Store(locationId, &emptyMap)

			logger.Debugf("NormalEvent end")
			return
		}
		allMap := allValue.(*sync.Map)
		node.Range(func(key, value interface{}) bool {
			logger.Errorf("NormalEvent values:", key, value)
			allMap.Delete(key)
			return true
		})
	} else {
		if clean {
			if ok {
				// 清盘时发送清事件消息
				node := value.(*sync.Map)
				events := this.GetMapEvents(node)
				if events != nil && len(events) > 0 {
					response := &proto.S2CBroadRandEvent{}
					response.LocationID = locationId
					response.EventInfo = events

					packet := &utils.Packet{}
					packet.Initialize(proto.MSG_BROAD_CLEAN_EVENT)
					packet.WriteData(response)

					logger.Debugf("NormalEvent BROAD_CLEAN_EVENT %+V %+V", locationId, events)
					Sched.BroadCastMsgToPlatform(int32(locationId), "", packet, proto.PLATFORM_ANDROID)
				}
			}
		}
		var emptyMap sync.Map
		G_CityRandMap.Store(locationId, &emptyMap)
	}
	logger.Debugf("NormalEvent end")
}

func (this *NormalEvent) GetMapEvents(nodeMap *sync.Map) map[int]*proto.EventData {
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

func (this *NormalEvent) EventProcess() {
	G_NormalEventTimer.Reset(time.Duration(G_NormalEventTimerReset) * time.Second)
	// logger.Debugf("NormalEvent EventProcess ", time.Now(), G_NormalEventTimerReset)
	if !G_NormalEventRunning {
		return
	}
	this.brocastEvent()
}

func (this *NormalEvent) Start() {
	G_NormalEventTimerReset = G_BaseCfg.EventInterval
	G_NormalEventTimer.Reset(time.Duration(1) * time.Second)

	G_NormalEventRunning = true
}

func (this *NormalEvent) Stop() {
	// if G_NormalEventRunning {
	G_NormalEventRunning = false
	G_NormalEventTimer.Stop()
	// }
}

//根据指定的源map随机取指定数量的元素,并返回取的源map的元素集合
func (this *NormalEvent) getRandPositionMap(locationId int, nodeRand *sync.Map) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData, 0)

	if len(G_BaseCfg.EventRate) == 0 {
		return ret
	}

	value, ok := G_CityAllMap.Load(locationId)
	if !ok {
		return ret
	}

	nodeAll := value.(*sync.Map)
	if utils.GetMapLen(nodeAll) == 0 {
		return ret
	}

	eventLen := utils.GetMapLen(nodeRand)
	generate_num := G_BaseCfg.EventNum - eventLen
	if generate_num <= 0 {
		return ret
	}

	leftMap := getLeftMap(locationId)
	if generate_num > len(leftMap) {
		generate_num = len(leftMap)
	}

	for i := 1; i <= generate_num; i++ {
		idx := rand.Intn(len(leftMap)) + 1
		positionNo := leftMap[idx]

		var node interface{}
		if node, ok = nodeAll.Load(positionNo); !ok {
			continue
		}
		value := node.(proto.EventNode)

		event_idx := rand.Intn(len(G_BaseCfg.EventRate))
		event_data := new(proto.EventData)
		event_data.ActivityType = proto.ACTIVITY_TYPE_NOMAL
		event_data.Type = value.Type
		event_data.X = value.X
		event_data.Y = value.Y
		event_data.EventId = getEventId(event_idx)
		ret[positionNo] = event_data
		nodeAll.Delete(positionNo)
		// G_CityAllMap.Store(locationId, nodeAll)
		// nodeRand.Store(positionNo, event_data)
		// G_CityRandMap.Store(locationId, nodeRand)
		leftMap = getLeftMap(locationId)
		//logger.Debugf("generate,positionNo=%+V event=%+V", positionNo, event_data)
	}
	return ret
}

func (s *CSession) HandleQueryShop(requestMsg *utils.Packet) {
	logger.Debugf("HandleQueryShop in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_QUERY_SHOP_RSP)

	responseMessage := &proto.S2CQueryShop{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SQueryShop{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	data, ret := db_service.WorldMapIns.QueryShop(int(s.LocateId), msg.KeyWord)
	if ret != nil {
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Data = make([]proto.ShopItem, 0)
	for _, value := range data {
		var item proto.ShopItem
		item.Desc = value.Desc
		item.Id = value.Id
		item.BuildingName = value.BuildingName
		item.ImageUrl = value.ImageUrl
		item.PassportAviable = value.PassportAviable
		item.PositionX = value.PositionX
		item.PositionY = value.PositionY
		item.ShopName = value.ShopName
		item.SmallType = value.SmallType
		item.H5Url = value.H5Url
		item.WebUrl = value.WebUrl
		responseMessage.Data = append(responseMessage.Data, item)
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleQueryShop end")

}
