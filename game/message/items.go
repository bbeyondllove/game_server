package message

import (
	"encoding/json"
	kk_core "game_server/core"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"game_server/core/logger"

	"github.com/shopspring/decimal"
)

var G_ItemList sync.Map

func GetItemList() bool {
	itemInfo, err := db_service.ItemIns.GetAllData()
	if err != nil {
		return false
	}

	for _, value := range itemInfo {
		node := &proto.ProductItem{
			ItemId:    value.Id,
			ItemType:  value.ItemType,
			ItemName:  value.ItemName,
			IsBind:    value.IsBind,
			IsGift:    value.IsGift,
			Quality:   value.Quality,
			IsPile:    value.IsPile,
			GetFrom:   value.GetFrom,
			UseJump:   value.UseJump,
			Price:     decimal.NewFromFloat32(value.Price),
			Recommend: value.Recommend,
			Desc:      value.Desc,
			Attr1:     value.Attr1,
			ImgUrl:    value.ImgUrl,
		}
		G_ItemList.Store(node.ItemId, node)
	}

	return true
}

//获取商品列表
func (s *CSession) HandleGetItemsList(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetItemsList in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_ITEMS_LIST_RSP)
	responseMessage := &proto.S2CItemList{}

	msg := &proto.C2SGetItemList{}
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

	node := make(map[int][]*proto.ProductItem, 0)
	itemLen := utils.GetMapLen(&G_ItemList)
	if itemLen == 0 {
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		responseMessage.ProductItemList = node
		rsp.WriteData(responseMessage)
		logger.Errorf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		return
	}

	itemIdAry := make([]int, 0)
	G_ItemList.Range(func(key interface{}, value interface{}) bool {
		item := value.(*proto.ProductItem)
		if item.IsBind == 1 {
			if _, ok := node[item.ItemType]; !ok {
				node[item.ItemType] = make([]*proto.ProductItem, 0)
			}

			itemIdAry = append(itemIdAry, item.ItemId)
		}
		return true
	})

	sort.Ints(itemIdAry)
	for _, subv := range itemIdAry {
		item, ok := G_ItemList.Load(subv)
		if !ok {
			continue
		}
		nodedata := item.(*proto.ProductItem)
		node[nodedata.ItemType] = append(node[nodedata.ItemType], nodedata)
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	responseMessage.ProductItemList = node

	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleGetItemsList end")
	return
}

func UpdateData(userId string, user_info map[string]string, key string, itemInfo *proto.ProductItem, itemNum int) error {
	logger.Debugf("UpdateData in request:%+v,%+v,%+v,%+v,%+v", userId, user_info, itemInfo, itemNum)
	var err error
	userItem := user_info[key]

	itemMap := make(map[int]map[int]int, 0)
	err = json.Unmarshal([]byte(userItem), &itemMap)
	if err != nil {
		return err
	}

	logger.Debugf("UpdateData Unmarshal:", itemMap)
	bFlag := false
	if _, ok := itemMap[itemInfo.ItemType]; ok {
		if _, subok := itemMap[itemInfo.ItemType][itemInfo.ItemId]; subok {
			itemMap[itemInfo.ItemType][itemInfo.ItemId] = itemMap[itemInfo.ItemType][itemInfo.ItemId] + itemNum
			bFlag = true
		} else {
			itemMap[itemInfo.ItemType][itemInfo.ItemId] = itemNum
		}
	} else {
		itemMap[itemInfo.ItemType] = make(map[int]int, 0)
		itemMap[itemInfo.ItemType][itemInfo.ItemId] = itemNum
	}
	logger.Debugf("UpdateData add map:%+V", itemMap)

	if bFlag {
		data_map := make(map[string]interface{})
		data_map["update_time"] = time.Now()
		data_map["item_num"] = itemMap[itemInfo.ItemType][itemInfo.ItemId]
		_, err = db_service.UserKnapsackIns.UpdateData(userId, itemInfo.ItemId, data_map)
	} else {
		userdata := model.UserKnapsack{
			UserId:   userId,
			ItemType: itemInfo.ItemType,
			ItemId:   itemInfo.ItemId,
			ItemNum:  itemNum,
		}
		_, err = db_service.UserKnapsackIns.Add(&userdata)
	}
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(itemMap)
	user_info[key] = string(buf)
	db.RedisGame.HSet(userId, key, buf)
	return err
}

//购买商品
func (s *CSession) HandleBuyItem(requestMsg *utils.Packet) {
	logger.Debugf("HandleBuyItem in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_BUY_ITEM_RSP)

	responseMessage := &proto.S2CBuyItem{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SBuyItem{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, user_info := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	logger.Debugf("in request:", msg)
	// logger.Debugf("user_info=", user_info)
	kk_core.PushMysql(func() {
		value, ok := G_ItemList.Load(msg.ItemId)
		if !ok {
			logger.Errorf("HandleBuyItem G_ItemList.Load failed(), itemid=", msg.ItemId)
			kk_core.PushWorld(func() {
				rsp.WriteData(responseMessage)
				s.sendPacket(rsp)
			})
			return
		}

		//会员购买限制
		//为了测试功能先屏蔽限制处理
		/*
			if value, ok := user_info["user_type"]; !ok || value != "1" {
				responseMessage.Code = errcode.ERROR_NOT_MEMBER
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				logger.Errorf("HandleBuyItem eror,非会员不能购买:", s.UserId)
				SendPacket(s.conn, responseMessage)
				return
			}
		*/
		itemInfo := value.(*proto.ProductItem)
		cdtvalue, _ := decimal.NewFromString(user_info["cdt"])

		//判断余额是否够
		if cdtvalue.Cmp(itemInfo.Price) < 0 {
			responseMessage.Code = errcode.ERROR_NOT_ENOUGH_MONEY
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			logger.Errorf("HandleBuyItem not enough money:", cdtvalue)
			kk_core.PushWorld(func() {
				rsp.WriteData(responseMessage)
				s.sendPacket(rsp)
			})
			return
		}

		if itemInfo.ItemType == proto.ITEM_LOCK { //角色解锁卡
			availableRoles := strings.Split(user_info["available_roles"], "|")
			if !utils.IsExistInArrs(itemInfo.Attr1, availableRoles) {
				availableRoles = append(availableRoles, itemInfo.Attr1)
			}

			rolesStr := strings.Join(availableRoles, "|")
			dataMap := make(map[string]interface{}, 0)
			dataMap["available_roles"] = rolesStr
			dataMap["update_time"] = time.Now()
			_, err = db_service.UpdateFields(db_service.UserTable, "user_id", payLoad.UserId, dataMap)
			if err != nil {
				logger.Errorf("UpdateFields %+v failed err=%+v", db_service.UserTable, err.Error())
				kk_core.PushWorld(func() {
					responseMessage.Code = errcode.ERROR_MYSQL
					responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
					rsp.WriteData(responseMessage)
					s.sendPacket(rsp)
				})
				return
			}

			db.RedisGame.HMSet(payLoad.UserId, dataMap)
		}

		price, _ := itemInfo.Price.Float64()
		code, lastCdt := db_service.NewCdt().UpdateUserCdt(payLoad.UserId, -float32(price), proto.MSG_BUY_ITEM)
		if code != errcode.MSG_SUCCESS {
			responseMessage.Code = errcode.ERROR_MYSQL
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			s.sendPacket(rsp)
			return
		}
		award := new(proto.AwardInfo)
		award.ItemId = itemInfo.ItemId
		award.ItemNum = 1
		award.ImgUrl = itemInfo.ImgUrl
		award.ItemName = itemInfo.ItemName
		award.Desc = itemInfo.Desc
		_, _, _ = sendAward(proto.MSG_FINISH_EVENT, payLoad.UserId, user_info, award)

		items := proto.ProductItem{
			ItemId:    itemInfo.ItemId,
			ItemType:  itemInfo.ItemType,
			ItemName:  itemInfo.ItemName,
			IsBind:    itemInfo.IsBind,
			Quality:   itemInfo.Quality,
			IsPile:    itemInfo.IsPile,
			GetFrom:   itemInfo.GetFrom,
			UseJump:   itemInfo.UseJump,
			Price:     itemInfo.Price,
			Recommend: itemInfo.Recommend,
			Desc:      itemInfo.Desc,
			ImgUrl:    itemInfo.ImgUrl,
		}
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = ""
		responseMessage.ItemId = items
		responseMessage.Money = decimal.NewFromFloat32(lastCdt)
		rsp.WriteData(responseMessage)
		logger.Debugf(payLoad.UserId + ":" + string(rsp.Bytes()))
		s.sendPacket(rsp)
		logger.Debugf("HandleBuyItem end")
	})

	return
}

func GetBindItemId() []string {
	ret := make([]string, 0)
	G_ItemList.Range(func(key interface{}, value interface{}) bool {
		item := value.(*proto.ProductItem)
		if item.Attr1 != "0" {
			ret = append(ret, item.Attr1)
		}
		return true
	})

	return ret
}

func GetItemId(roleId int) int {
	ret := 0
	G_ItemList.Range(func(key interface{}, value interface{}) bool {
		item := value.(*proto.ProductItem)
		if item.Attr1 == strconv.Itoa(roleId) {
			ret = item.ItemId
		}
		return true
	})

	return ret
}

func GetRoleId(itemId int) string {
	ret := ""
	G_ItemList.Range(func(key interface{}, value interface{}) bool {
		item := value.(*proto.ProductItem)
		if item.ItemId == itemId {
			ret = item.Attr1
		}
		return true
	})

	return ret
}
