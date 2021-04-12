package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/proto"
	"sort"
	"strconv"
	"time"
)

func GetUserItem(userInfo map[string]string, key string) (map[int]map[int]int, map[int][]proto.ItemInfo) {
	userItem := userInfo[key]
	itemMap := make(map[int]map[int]int, 0)
	itemInfos := make(map[int][]proto.ItemInfo, 0)
	err := json.Unmarshal([]byte(userItem), &itemMap)
	if err != nil {
		return itemMap, itemInfos
	}

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

	return itemMap, itemInfos
}

func IsUserHaveItem(userInfo map[string]string, itemId int, key string) (bool, *proto.ProductItem, map[int]map[int]int) {
	var item *proto.ProductItem
	node, ok := G_ItemList.Load(itemId)
	if !ok {
		return false, nil, nil
	}

	item = node.(*proto.ProductItem)
	itemMap, _ := GetUserItem(userInfo, key)

	bFlag := false
	if _, ok = itemMap[item.ItemType]; ok {
		if _, subok := itemMap[item.ItemType][itemId]; subok {
			if itemMap[item.ItemType][itemId] > 0 {
				bFlag = true
			}
		}
	}

	return bFlag, item, itemMap
}

func UseItem(userId string, userInfo map[string]string, itemId int) (bool, int32) {
	bFlag, item, itemMap := IsUserHaveItem(userInfo, itemId, "item_info")
	if !bFlag {
		return false, errcode.ERROR_NO_SIGNIN_LOST_CARD
	}

	data_map := make(map[string]interface{}, 0)
	data_map["update_time"] = time.Now()
	bDelete := false
	if itemMap[item.ItemType][itemId] == 1 {
		delete(itemMap[item.ItemType], itemId)
		bDelete = true
	} else {
		itemMap[item.ItemType][itemId]--
	}
	userItemInfo, _ := json.Marshal(itemMap)
	data_map["item_info"] = userItemInfo
	_, err := db.RedisGame.HMSet(userId, data_map).Result()
	if err != nil {
		return false, errcode.ERROR_SYSTEM
	}

	//
	if bDelete {
		_, err = db_service.UserKnapsackIns.Delete(userId, itemId)
		if err != nil {
			itemMap[item.ItemType][itemId] = 1
		}
	} else {
		data_map = make(map[string]interface{}, 0)
		data_map["update_time"] = time.Now()
		data_map["item_num"] = itemMap[item.ItemType][itemId]

		_, err = db_service.UserKnapsackIns.UpdateData(userId, itemId, data_map)
		if err != nil {
			itemMap[item.ItemType][itemId]++
		}
	}

	if err != nil {
		userItemInfo, _ := json.Marshal(itemMap)
		data_map["item_info"] = userItemInfo
		_, err = db.RedisGame.HMSet(userId, data_map).Result()
		return false, errcode.ERROR_SYSTEM
	}

	userInfo["item_info"] = string(userItemInfo)
	return true, errcode.MSG_SUCCESS
}

func GetInviteUsers(userId string, level, page, size, vipLevel int) (error, *proto.S2CInviteUserList) {
	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = proto.SysType
	request_map["inviterId"] = userId
	request_map["invitationLevel"] = level
	request_map["page"] = page
	request_map["size"] = size
	request_map["vipLevel"] = vipLevel
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	ret := &proto.S2CInviteUserList{}
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetInviteUsersUrl
	msgdata, err := utils.HttpGet(url, request_map)
	if err != nil {
		return err, ret
	}

	err = json.Unmarshal(msgdata, ret)
	return err, ret
}
