package message

import (
	"encoding/json"
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/timer"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/message/activity/activity_roles"
	"game_server/game/message/activity/double_year"
	"game_server/game/message/statistical"
	"game_server/game/model"
	"game_server/game/proto"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// DoubleYearEvent 双旦活动事件struct.
type DoubleYearEvent struct {
	rankInstance  *double_year.RankList
	cdtInstance   *double_year.TradeCdt
	cdtDbInstance *db_service.Cdt
}

var (
	G_ChristmaAllMap    sync.Map //所有地图，key1城市编号,key2坐标编号 value "x,y"坐标
	G_ChristmaRandMap   sync.Map //随机地图，key1城市编号,key2坐标编号 value 事件数据
	G_ChristmaEventConf sync.Map //key 事件序号，value事件信息

)

const (
	CHRISTMA_MAP          = "christma_map"
	CHRISTMA_SUIPIAN      = "christma_suipian:"
	TREE_SWEET_KEY        = "item_info"
	PUSH_CUMULATIVE_POPUP = "push_cumulative_popup:"
)

func init() {
	G_DoubleYearEvent.rankInstance = double_year.NewRankList()
	G_DoubleYearEvent.cdtInstance = new(double_year.TradeCdt)
	G_DoubleYearEvent.cdtDbInstance = db_service.NewCdt()

}

func (this *DoubleYearEvent) loadItem() {
	userItem, err := db_service.DoubleYearIns.GetAllData()
	if err != nil || len(userItem) == 0 {
		return
	}

	for _, v := range userItem {
		var itemValue interface{}
		ok := false
		if itemValue, ok = G_ItemList.Load(v.ItemId); !ok {
			continue
		}

		item := itemValue.(*proto.ProductItem)
		this.setSuiPian(v.UserId, item, v.ItemNum, true)
	}
}

func (this *DoubleYearEvent) loadMap() {
	dataMap := db.RedisMgr.HGetAll(CHRISTMA_MAP)
	if dataMap == nil || len(dataMap) == 0 {
		return
	}

	for k, v := range dataMap {
		locatinId, _ := strconv.Atoi(k)
		node := make(map[int]*proto.EventData)
		err := json.Unmarshal([]byte(v), &node)
		if err != nil {
			continue
		}
		var randMap sync.Map
		for key, value := range node {
			randMap.Store(key, value)
		}
		G_ChristmaRandMap.Store(locatinId, &randMap)
	}
	this.loadItem()
}

//func (this *DoubleYearEvent) getMap() map[int]map[int]*proto.EventData {
//	ret := make(map[int]map[int]*proto.EventData)
//	G_ChristmaRandMap.Range(func(k interface{}, v interface{}) bool {
//		lid := k.(int)
//		ret[lid] = make(map[int]*proto.EventData, 0)
//		data := v.(*sync.Map)
//		data.Range(func(key interface{}, value interface{}) bool {
//			positionId := key.(int)
//			ret[lid][positionId] = value.(*proto.EventData)
//			return true
//		})
//		return true
//	})
//	return ret
//}

func (this *DoubleYearEvent) getMap(locationId int) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData)
	value, ok := G_ChristmaRandMap.Load(locationId)
	if ok {
		nodeRand := value.(*sync.Map)
		nodeRand.Range(func(key interface{}, v interface{}) bool {
			positionId := key.(int)
			ret[positionId] = v.(*proto.EventData)
			return true
		})
		return ret
	}
	return ret
}

func (this *DoubleYearEvent) getActivityState() int {
	startDate := utils.Str2Time(base.Setting.Springfestival.ActivityStartDatetime)
	endDate := utils.Str2Time(base.Setting.Springfestival.ActivityEndDatetime)
	curDate := time.Now()
	if curDate.Sub(startDate) < 0 {
		return proto.ACTIVITY_NOT_START
	}

	if curDate.Sub(endDate) >= 0 {
		return proto.ACTIVITY_END
	}

	return proto.ACTIVITY_START
}

func (this *DoubleYearEvent) Start() {
	timer.SetTimer("startEvent", 1, this.startEvent, nil)
	go timer.Run()
}

func (this *DoubleYearEvent) sendAward(v interface{}) {
	userItem, err := db_service.DoubleYearIns.GetAllData()
	if err != nil || len(userItem) == 0 {
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["is_trade"] = 1
	dataMap["update_time"] = time.Now()

	for _, v := range userItem {
		if v.ItemNum < base.Setting.Doubleyear.ChristmasPatchTradeNum {
			continue
		}

		_, err = db.RedisGame.HSet(v.UserId, "is_trade", 1).Result()
		if err != nil {
			continue
		}

		_, err = db_service.DoubleYearIns.UpdateData(v.UserId, v.ItemId, dataMap)
		if err != nil {
			db.RedisGame.HSet(v.UserId, "is_trade", 0)
			continue
		}

		email := make(map[string]interface{})
		email["userId"] = v.UserId
		email["emailType"] = 2
		email["emailTitle"] = base.Setting.Doubleyear.SantaEmailTitle
		email["emailContent"] = base.Setting.Doubleyear.SantaEmailContent
		// 角色ID
		rolesId := base.Setting.Springfestival.GoldenCoupleRole
		//获取角色card_id
		emaiPrize := make([]map[string]interface{}, 0)
		for _, roleId := range rolesId {
			activit := activity_roles.NewActivityRoles()
			cardId := activit.GetItemIdByRoleId(roleId)
			node, ok := G_ItemList.Load(cardId)
			if !ok {
				logger.Errorf("ItemId error %+V", cardId)
				return
			}
			itemInfo := node.(*proto.ProductItem)
			prizeList := map[string]interface{}{
				"prizeType": 1,                 //奖励类型
				"prizeId":   itemInfo.ItemId,   // 奖励id
				"prizeName": itemInfo.ItemName, //奖励名称
				"prizeNum":  1,                 //奖励数量
				"prizeImg":  itemInfo.ImgUrl,   // 图片
			}
			emaiPrize = append(emaiPrize, prizeList)
			// 只添加其中一个
			break
		}
		email["prizeList"] = emaiPrize
		_, err := db_service.EmailLogicIns.AddEmail(email)
		logger.Infof("DoubleYear sendAward return %+V %+V", v.UserId, err)

	}
}

func (this *DoubleYearEvent) startEvent(v interface{}) {
	//logger.Debugf("startEvent in")
	status := this.getActivityState()

	timer.Delete(timer.TimerMap["startEvent"])
	if status == proto.ACTIVITY_NOT_START {
		timer.SetTimer("startEvent", 1, this.startEvent, v)
	} else {
		timer.SetTimer("BrocastEvent", 1, this.BrocastEvent, v)
	}
	//logger.Debugf("startEvent end")

}

//初始化
func (this *DoubleYearEvent) EventInit() {
	logger.Debugf("DoubleYearEventInit in")

	G_ChristmaAllMap.Range(func(key interface{}, value interface{}) bool {
		G_ChristmaAllMap.Delete(key)
		return true
	})
	G_ChristmaRandMap.Range(func(key interface{}, value interface{}) bool {
		G_ChristmaRandMap.Delete(key)
		return true
	})

	logger.Debugf("DoubleYearEventInit end")
}

// 时间判断
func (this *DoubleYearEvent) isInTime(eventRate *proto.EventRate) bool {
	status := this.getActivityState()
	if status != proto.ACTIVITY_START {
		return false
	}

	if eventRate.StartTime == "" && eventRate.StartDate == "" {
		return true
	}

	currentTime := time.Now()
	// 判断天
	if currentTime.Unix() < utils.Str2Time(eventRate.StartDate+" "+"00:00:00").Unix() || currentTime.Unix() >= utils.Str2Time(eventRate.EndDate+" "+"00:00:00").Unix() {
		return false
	}
	//判断小时
	currentDate := utils.GetCurDay()
	if currentTime.Unix() < utils.Str2Time(currentDate+" "+eventRate.StartTime).Unix() || currentTime.Unix() >= utils.Str2Time(currentDate+" "+eventRate.EndTime).Unix() {
		return false
	}

	return true
}

//存储节点
func (this *DoubleYearEvent) EventStore(locatinId int, data sync.Map) {
	logger.Debugf("doubleYearStore in")

	G_ChristmaAllMap.Store(locatinId, &data)
	var emptyMap sync.Map
	G_ChristmaRandMap.Store(locatinId, &emptyMap)

	logger.Debugf("doubleYearStore end")
}

func (this *DoubleYearEvent) getItemCount(itemId int) int {
	key := "christma_map_" + strconv.Itoa(itemId)
	field := "count"
	valueRet, err := db.RedisGame.HGet(key, field).Result()
	if err == nil && valueRet != "" {
		value, _ := strconv.Atoi(valueRet)
		return value
	}
	return 0
}

func (this *DoubleYearEvent) setItemCount(itemId int, num int) {
	key := "christma_map_" + strconv.Itoa(itemId)
	field := "count"
	utils.SetKeyValue(key, field, int64(num), true, utils.ITEM_DAY)

}
func (this *DoubleYearEvent) generateItem(locationId int, nodeAll *sync.Map, nodeRand *sync.Map, eventRate *proto.EventRate, limitNum int) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData)
	if !this.isInTime(eventRate) {
		return ret
	}
	/*	count := this.getItemCount(eventRate.ItemId)
		if count >= base.Setting.Doubleyear.ChristmasUnitMaxNum {
			//logger.Debugf("generateItem cout is to max", count)
			return ret
		}
	*/
	generateNum := eventRate.UnitNum
	/*	canRandNum := base.Setting.Doubleyear.ChristmasUnitMaxNum - count
		if canRandNum < generateNum {
			generateNum = canRandNum
		}
	*/
	// 判断事件生成的长度
	generate_num := eventRate.LimitNum - limitNum
	if generate_num <= 0 {
		return ret
	}
	if generate_num < generateNum {
		generateNum = generate_num
	}
	for i := 0; i < generateNum; i++ {
		leftMap := this.getLeftMap(locationId)
		lenLeft := len(leftMap)

		if lenLeft == 0 {
			return ret
		}

		idx := rand.Intn(lenLeft) + 1
		positionNo := leftMap[idx]
		var node interface{}
		var ok bool
		if node, ok = nodeAll.Load(positionNo); !ok {
			return ret
		}
		eventNode := node.(proto.EventNode)
		event_data := new(proto.EventData)
		event_data.ActivityType = eventNode.ActivityType
		event_data.Type = eventNode.Type
		event_data.X = eventNode.X
		event_data.Y = eventNode.Y
		event_data.EventId = eventRate.Id

		nodeAll.Delete(positionNo)
		// G_ChristmaAllMap.Store(locationId, nodeAll)
		nodeRand.Store(positionNo, event_data)
		// G_ChristmaRandMap.Store(locationId, nodeRand)
		this.setItemCount(eventRate.ItemId, int(eventRate.ItemNum))
		ret[positionNo] = event_data
	}
	return ret

}

//产生并广播随机事件
func (this *DoubleYearEvent) BrocastEvent(v interface{}) {
	//logger.Debugf("DoubleYearEvent brocastEvent in")
	status := this.getActivityState()
	if status == proto.ACTIVITY_END {
		timer.Delete(timer.TimerMap["BrocastEvent"])
		timer.Stop()
		this.EventInit()
		// 活动结束删除掉随机事件缓存
		db.RedisGame.Del(CHRISTMA_MAP)

		go this.sendAward(v)
		//广播活动结束
		rsp := &utils.Packet{}
		rsp.Initialize(proto.MSG_BROAD_ACTIVITY_END)
		responseMessage := &proto.S2CActivityEnd{}
		responseMessage.ActivityType = proto.ACTIVITY_TYPE_DOUBLE_YEAR
		rsp.WriteData(responseMessage)
		Sched.SendToAllUser(rsp)
		//广播活动状态
		go PushActivityStatus(true, "")
		logger.Infof("DoubleYear activity is end")
		return
	}

	timer.Delete(timer.TimerMap["BrocastEvent"])
	interval := uint32(base.Setting.Doubleyear.ChristmasEventInterval)

	for _, v := range G_BaseCfg.LocationId {
		eventData := this.getRandPositionMap(v)
		if len(eventData) == 0 {
			continue
		}

		dataMap := this.getMap(v)
		buf, _ := json.Marshal(dataMap)
		utils.SetKeyValue(CHRISTMA_MAP, strconv.Itoa(v), buf, false, utils.ITEM_DAY)

		//logger.Debugf("brocastEvent retEvent=%+V %+V", v, eventData)
		rsp := &utils.Packet{}
		rsp.Initialize(proto.MSG_BROAD_RAND_EVENT)
		responseMessage := &proto.S2CBroadRandEvent{}
		responseMessage.LocationID = v
		responseMessage.EventInfo = eventData
		rsp.WriteData(responseMessage)
		logger.Debugf("DoubleYearEvent brocastEvent %+V %+V", v, eventData)

		Sched.BroadCastMsg(int32(v), "", rsp)
	}
	timer.SetTimer("BrocastEvent", interval, this.BrocastEvent, v)

	//logger.Debugf("DoubleYearEvent brocastEvent end")
}

func (this *DoubleYearEvent) getLeftMap(locationId int) map[int]int {
	ret := make(map[int]int, 0)
	i := 1
	if value, ok := G_ChristmaAllMap.Load(locationId); ok {
		node := value.(*sync.Map)
		node.Range(func(key interface{}, value interface{}) bool {
			ret[i] = key.(int)
			i++
			return true
		})
	}

	return ret
}

//根据指定的源map随机取指定数量的元素,并返回取的源map的元素集合
func (this *DoubleYearEvent) getRandPositionMap(locationId int) map[int]*proto.EventData {
	ret := make(map[int]*proto.EventData, 0)

	if len(G_BaseCfg.DoubleEventRate) == 0 {
		return ret
	}

	value, ok := G_ChristmaAllMap.Load(locationId)
	if !ok {
		return ret
	}

	nodeAll := value.(*sync.Map)
	if utils.GetMapLen(nodeAll) == 0 {
		return ret
	}

	value, ok = G_ChristmaRandMap.Load(locationId)
	if !ok {
		return ret
	}

	nodeRand := value.(*sync.Map)
	mapLen := utils.GetMapItemIdLen(nodeRand) // 获取生成事件的数量
	G_ChristmaEventConf.Range(func(key interface{}, value interface{}) bool {
		nodeEvent := value.(proto.EventRate)
		eventData := this.generateItem(locationId, nodeAll, nodeRand, &nodeEvent, mapLen[nodeEvent.ItemId])
		if len(eventData) > 0 {
			for k, v := range eventData {
				ret[k] = v
				//logger.Debugf("generate,positionNo=%+V event=%+V", k, v)
			}
		}
		return true
	})

	return ret
}

//点击事件
func (this *DoubleYearEvent) DoubleYearFinishEvent(userId string, userInfo map[string]string, eventRate *proto.EventRate) (bool, map[int][]*proto.AwardItem) {
	logger.Debugf("DoubleYearFinishEvent in request:", userId, eventRate)
	ret := make(map[int][]*proto.AwardItem)
	var itemInfo *proto.ProductItem
	node, ok := G_ItemList.Load(eventRate.ItemId)
	if !ok {
		logger.Errorf("ItemId error %+V", eventRate)
		return false, ret
	}
	itemInfo = node.(*proto.ProductItem)

	if eventRate.LimitType >= 3 {
		this.UpdateItem(userId, userInfo, itemInfo, int(eventRate.ItemNum), eventRate.LimitType)
	} else {
		this.rankInstance.UpdateProp(userId, eventRate.ItemId, int(eventRate.ItemNum))
	}

	ret[itemInfo.ItemType] = make([]*proto.AwardItem, 0)
	award := new(proto.AwardItem)
	award.ItemId = eventRate.ItemId
	award.ItemNum = eventRate.ItemNum
	award.ImgUrl = itemInfo.ImgUrl
	award.ItemName = itemInfo.ItemName
	award.Desc = itemInfo.Desc
	ret[itemInfo.ItemType] = append(ret[itemInfo.ItemType], award)
	logger.Debugf("DoubleYearFinishEvent end")
	return true, ret
}

func (this *DoubleYearEvent) getSuiPian(userId string) map[string]interface{} {
	key := CHRISTMA_SUIPIAN + userId
	dataMap := db.RedisMgr.HGetAll(key)
	ret := make(map[string]interface{})
	if dataMap == nil {
		return ret
	}

	item_type, _ := strconv.Atoi(dataMap["item_type"])
	item_id, _ := strconv.Atoi(dataMap["item_id"])
	num, _ := strconv.Atoi(dataMap["num"])
	is_trade, _ := strconv.Atoi(dataMap["is_trade"])
	ret["item_type"] = item_type
	ret["item_id"] = item_id
	ret["num"] = num
	ret["user_id"] = dataMap["user_id"]
	ret["is_trade"] = is_trade
	logger.Debugf("DoubleYearFinishEvent getSuiPian dataMap=%+V,ret=%+V", dataMap, ret)

	return ret
}

func (this *DoubleYearEvent) setSuiPian(userId string, itemInfo *proto.ProductItem, itemNum int, bReset bool) (error, int) {
	key := CHRISTMA_SUIPIAN + userId
	oldMap := this.getSuiPian(userId)
	dataMap := make(map[string]interface{})
	dataMap["user_id"] = userId
	dataMap["item_type"] = itemInfo.ItemType
	dataMap["item_id"] = itemInfo.ItemId
	dataMap["num"] = itemNum
	oldNum := oldMap["num"].(int)

	if oldNum == 0 || bReset {
		_, err := db.RedisGame.HMSet(key, dataMap).Result()
		if err != nil {
			return err, 0
		}
		// 新春活动结束时间
		//utils.SetKeyTiemExp(key, base.Setting.Doubleyear.PropEndDate+" "+base.Setting.Doubleyear.PropEndTime)
		utils.SetKeyTiemExp(key, base.Setting.Springfestival.ActivityEndDatetime)
		userdata := model.DoubleYearUserItem{
			UserId:   userId,
			ItemType: itemInfo.ItemType,
			ItemId:   itemInfo.ItemId,
			ItemNum:  itemNum,
		}
		_, err = db_service.DoubleYearIns.Add(&userdata)
		if err != nil {
			return err, 0
		}

		return nil, itemNum
	}
	dataMap["num"] = oldNum + itemNum
	logger.Debugf("DoubleYearFinishEvent setSuiPian oldNum=%+v,oldNum=%+v,map=%+V", oldNum, itemNum, dataMap)
	_, err := db.RedisGame.HMSet(key, dataMap).Result()
	if err != nil {
		return err, 0
	}

	db_map := make(map[string]interface{})
	db_map["update_time"] = time.Now()
	db_map["item_num"] = dataMap["num"]
	_, err = db_service.DoubleYearIns.UpdateData(userId, itemInfo.ItemId, db_map)
	return err, oldNum + itemNum
}

func (this *DoubleYearEvent) getTreeAndSweet(userId string) map[int]map[int]int {
	key := CHRISTMA_MAP + ":" + userId

	itemMap := make(map[int]map[int]int, 0)
	userItem, err := db.RedisGame.HGet(key, TREE_SWEET_KEY).Result()
	if err != nil && err.Error() != "redis: nil" {
		logger.Errorf("redis error", err)
		return itemMap
	}

	err = json.Unmarshal([]byte(userItem), &itemMap)
	if err != nil {
		return itemMap
	}

	return itemMap
}

func (this *DoubleYearEvent) setTreeAndSweet(userId string, itemInfo *proto.ProductItem, itemNum int) error {
	key := CHRISTMA_MAP + ":" + userId
	userItem := this.getTreeAndSweet(userId)

	if len(userItem) == 0 {
		userItem[itemInfo.ItemType] = make(map[int]int)
		userItem[itemInfo.ItemType][itemInfo.ItemId] = itemNum
	} else {
		if value, ok := userItem[itemInfo.ItemType][itemInfo.ItemId]; ok {
			userItem[itemInfo.ItemType][itemInfo.ItemId] = itemNum + value
		} else {
			userItem[itemInfo.ItemType][itemInfo.ItemId] = itemNum
		}
	}

	if userItem[itemInfo.ItemType][itemInfo.ItemId] == 0 {
		delete(userItem[itemInfo.ItemType], itemInfo.ItemId)
		if len(userItem[itemInfo.ItemType]) == 0 {
			delete(userItem, itemInfo.ItemType)
		}
	}

	buf, _ := json.Marshal(userItem)
	utils.SetKeyValue(key, TREE_SWEET_KEY, buf, false, utils.ITEM_DAY)
	return nil
}

func (this *DoubleYearEvent) UpdateItem(userId string, userInfo map[string]string, itemInfo *proto.ProductItem, itemNum, limitType int) error {
	logger.Debugf("DoubleYearEvent UpdateItem in request:%+v,%+v,%+v,%+v,%+v", userId, itemInfo, itemNum, limitType)
	var err error
	var num int

	if limitType > 3 {
		statistical.StatisticsDotIns.DoubleYearFuwa(userId, itemInfo.ItemId, itemNum)

		err, num = this.setSuiPian(userId, itemInfo, itemNum, false)
		if err == nil && num == base.Setting.Doubleyear.ChristmasPatchTradeNum {
			//广播跑马灯
			rsp := &utils.Packet{}
			rsp.Initialize(proto.MSG_BROAD_SANTA_CARD)
			responseMessage := &proto.S2CBroadSantaCard{}
			responseMessage.Content = strings.Replace(base.Setting.Doubleyear.SantaBroadContent, "%s", userInfo["nick_name"], 1)
			rsp.WriteData(responseMessage)
			Sched.SendToAllUser(rsp)

		}
	} else {
		err = this.setTreeAndSweet(userId, itemInfo, itemNum)
		if err == nil {
			this.rankInstance.UpdateProp(userId, itemInfo.ItemId, itemNum)
		}
	}

	return err
}

//用户双蛋圣诞树，糖果 （新春活动 幸运箭矢，幸福天使）
func (s *CSession) HandleGetSweetAndTree(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetSweetAndTree in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CGetSweetTree{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_SWEET_AND_TREE_RSP)

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	itemMap := G_DoubleYearEvent.getTreeAndSweet(payLoad.UserId)
	itemInfos := make(map[int][]proto.ItemInfo, 0)
	for k, v := range itemMap {
		itemInfos[k] = make([]proto.ItemInfo, 0)
		itemIdx := make([]int, 0)
		dataMap := make(map[int]proto.ItemInfo, 0)
		for key, value := range v {
			var item *proto.ProductItem
			ok := false
			var itemValue interface{}
			if itemValue, ok = G_ItemList.Load(key); !ok {
				continue
			}
			item = itemValue.(*proto.ProductItem)
			sex := 0
			roldid, _ := strconv.Atoi(item.Attr1)
			if roldid%2 == 0 {
				sex = 1
			}
			node := proto.ItemInfo{
				ItemId:   key,
				Num:      value,
				Desc:     item.Desc,
				Attr1:    item.Attr1,
				ItemName: item.ItemName,
				ImgUrl:   item.ImgUrl,
				Sex:      sex,
			}
			dataMap[key] = node
			itemIdx = append(itemIdx, key)
		}
		sort.Ints(itemIdx)
		for _, value := range itemIdx {
			itemInfos[k] = append(itemInfos[k], dataMap[value])
		}
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	userDayCdt, _ := G_DoubleYearEvent.cdtInstance.GetUserCdt(payLoad.UserId)
	totalDayCdt, _ := G_DoubleYearEvent.cdtInstance.GetTotalCdt()
	userAllCdt := G_DoubleYearEvent.cdtDbInstance.GetUserCdt(payLoad.UserId)

	responseMessage.SweetTree = itemInfos
	responseMessage.UserDayCdt = decimal.NewFromFloat32(userDayCdt)
	responseMessage.TotalDayCdt = decimal.NewFromFloat32(totalDayCdt)
	responseMessage.CurrentCdt = decimal.NewFromFloat32(userAllCdt)

	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleGetSweetAndTree end")
	return
}

//用户双蛋碎片
func (s *CSession) HandleGetPatch(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetPatch in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CGetPatch{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_PATCH_RSP)

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	currentNum := 0

	suiPianMap := G_DoubleYearEvent.getSuiPian(payLoad.UserId)
	if len(suiPianMap) > 0 {
		currentNum = suiPianMap["num"].(int)
	}

	responseMessage.CurrentNum = currentNum
	responseMessage.TradeNeedNum = base.Setting.Doubleyear.ChristmasPatchTradeNum
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleGetPatch end")
	return
}

//用户糖果圣诞树兑换cdt
func (s *CSession) HandleTradeCdt(requestMsg *utils.Packet) {
	logger.Debugf("HandleTradeCdt in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CTradeCdt{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SWEET_TREE_RSP)

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	dataMap := G_DoubleYearEvent.getTreeAndSweet(payLoad.UserId)
	if len(dataMap) == 0 {
		responseMessage.Code = errcode.ERROR_NOT_ENOUGH_ITEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

		logger.Errorf("HandleTradeCdt %+V", responseMessage.Message)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	sweetNum := 0
	treeNum := 0
	for _, v := range dataMap {
		if value, ok := v[proto.ITEM_SWEET]; ok {
			sweetNum = value
		}

		if value, ok := v[proto.ITEM_TREE]; ok {
			treeNum = value
		}
	}

	useSweetNum := 0
	useTreeNum := 0
	retMap, retErr := G_DoubleYearEvent.cdtInstance.TradeCdt(payLoad.UserId, sweetNum, treeNum)
	responseMessage.Code = int32(retErr)
	if responseMessage.Code == errcode.MSG_SUCCESS {
		useSweetNum = retMap["sweetUsed"].(int)
		useTreeNum = retMap["treeUsed"].(int)

		if useSweetNum > 0 {
			if itemValue, ok := G_ItemList.Load(proto.ITEM_SWEET); ok {
				node := itemValue.(*proto.ProductItem)
				G_DoubleYearEvent.setTreeAndSweet(payLoad.UserId, node, -useSweetNum)
			}
		}

		if useTreeNum > 0 {
			if itemValue, ok := G_ItemList.Load(proto.ITEM_TREE); ok {
				node := itemValue.(*proto.ProductItem)
				G_DoubleYearEvent.setTreeAndSweet(payLoad.UserId, node, -useTreeNum)
			}
		}
	}

	responseMessage.Message = double_year.StatusCodeMessage[retErr]
	responseMessage.UserId = payLoad.UserId
	changeCdt, _ := FormatFloat(float64(retMap["changeCdt"].(float32)), 4)
	responseMessage.ChangeCdt = float32(changeCdt)
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleTradeCdt end")
	// 兑换成功才推送，广播消息
	if responseMessage.Code == 0 {
		// 推送跑马灯 广播消息
		go PushRedEnvelopeMsg(payLoad.UserId, float32(changeCdt), 0)
		// todo 暂时分两调消息推送
		userDayCdt := retMap["userDayCdt"].(string)
		userDayTotalCdt, _ := strconv.ParseFloat(userDayCdt, 64)
		go PushCumulativePopup(payLoad.UserId, float32(userDayTotalCdt))
	}
	return
}

// 推送累计消息
func PushCumulativePopup(userId string, cdt float32) {
	var MsgMap = map[int]int{
		1:   0,
		5:   0,
		10:  0,
		20:  0,
		40:  0,
		60:  0,
		80:  0,
		100: 0,
	}
	key := PUSH_CUMULATIVE_POPUP + time.Now().Format("20060102") + ":" + userId
	result := db.RedisMgr.HGetAll(key)
	if result == nil || len(result) == 0 {
		pushDate := MsgMap
		for k, v := range pushDate {
			if cdt >= float32(k) && v == 0 {
				pushDate[k] = 1
				db.RedisMgr.HSet(key, strconv.Itoa(k), 1)
				go PushRedEnvelopeMsg(userId, float32(k), 1)
			} else {
				db.RedisMgr.HSet(key, strconv.Itoa(k), 0)
			}
		}
		db.RedisMgr.Expire(key, 3600*24) // 1个小时后过期
		return
	}

	for k, v := range result {
		kKey, _ := strconv.Atoi(k)
		kValue, _ := strconv.Atoi(v)
		if cdt >= float32(kKey) && kValue == 0 {
			db.RedisMgr.HSet(key, strconv.Itoa(kKey), 1)
			go PushRedEnvelopeMsg(userId, float32(kKey), 1)
		}
	}
}

func PushRedEnvelopeMsg(userId string, cdt float32, isTotal int) {
	//广播跑马灯
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_RED_ENVELOPE_RESP)
	responseMessage := &proto.S2CRedEnvelopeMsg{}
	responseMessage.Cdt = cdt
	// 获取用户昵称
	value, _ := db.RedisGame.HGet(userId, "nick_name").Result()
	responseMessage.UserName = value
	responseMessage.IsTotal = isTotal
	rsp.WriteData(responseMessage)
	Sched.SendToAllUser(rsp)
}

//获取活动状态
func (s *CSession) HandleGetDoubleYearStatus(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetDoubleYearStatus in request:", requestMsg.GetBuffer())

	responseMessage := &proto.S2CDoubleYearStatus{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_DOUBLE_YEAR_STATUS_RSP)

	msg := &proto.C2SGetActiveStatus{}
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

	//0:双旦;1:排行榜;
	if msg.StatusType == 1 {
		rankList := double_year.NewRankList()
		responseMessage.Status = rankList.GetRandListStatus()
	} else {
		responseMessage.Status = G_DoubleYearEvent.getActivityState()
	}

	responseMessage.StatusType = msg.StatusType
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)

	logger.Debugf("HandleGetDoubleYearStatus end")
	return
}

func (this *DoubleYearEvent) Test() {
	status := this.getActivityState()
	fmt.Println(status)
}
