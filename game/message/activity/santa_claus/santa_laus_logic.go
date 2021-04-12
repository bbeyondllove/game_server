package santa_claus

import (
	"encoding/json"
	"errors"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/model"
	"strconv"
	"strings"
	"time"
)

type SantaClaus struct {
}

var SantaClausLogic = NewSantaClaus()

func NewSantaClaus() *SantaClaus {
	return &SantaClaus{}
}

const SANTA_CLAUS_ROLES = "santa_claus:roles:" // 圣诞老人角色过期时间

//检查当前角色是否是圣诞老人
//func (this *SantaClaus) CheckCurrentRole(userId string) (bool, error) {
//	//获取当前用户的角色id
//	role_id, err := db.RedisMgr.GetRedisClient().HGet(userId, "role_id").Result()
//	if err != nil {
//		logger.Errorf("get user role_id get redis err:=", err.Error(), "userId = ", userId)
//		return false, err
//	}
//	roleId, _ := strconv.Atoi(role_id)
//	// 判断当前用户角色是不是圣诞老人
//	if roleId != base.Setting.Doubleyear.SantaClausRoleId {
//		return false, errors.New("It's not Santa Claus")
//	}
//	// 检查角色是否过期是否可用
//	isExpire, err := this.CheckUserRoleExpire(userId)
//	if isExpire == false {
//		// 移除用户角色
//		this.RemoveUserRole(userId)
//		// 通知用户已过期
//		//go this.NotifyUsersRoleIsExpire(userId)
//		return false, err
//	}
//	return isExpire, err
//}

//检查用户角色是否过期
func (this *SantaClaus) CheckUserRoleExpire(userId string) (bool, error) {
	//// 检查30天的时间是否已过期
	//isExpire, err := this.CheckEndExpire()
	//if isExpire == false {
	//	return isExpire, err
	//}
	//检查角色时间是否已过期
	rolesIsExpire, err := this.CheckSantaLausRoles(userId)
	if rolesIsExpire == false {
		return rolesIsExpire, err
	}
	return true, nil
}

//解锁卡有效期30天
//检查是否过期 false 已过期  true 可以选择
func (this *SantaClaus) CheckEndExpire() (bool, error) {
	currentTime := time.Now().Unix()
	loc, _ := time.LoadLocation("Local")                                                                             //获取时区
	expireTime, _ := time.ParseInLocation("2006-01-02 15:04:05", base.Setting.Doubleyear.SantaClausExpriteTime, loc) // 当前时间
	if currentTime > expireTime.Unix() {
		return false, errors.New("currentTime > expireTime")
	}
	return true, nil
}

//检查圣诞老人角色是否过期
func (this *SantaClaus) CheckSantaLausRoles(userId string) (bool, error) {
	key := SANTA_CLAUS_ROLES + userId
	expireTime, _ := db.RedisMgr.GetRedisClient().HGet(key, "ExpireTime").Result()
	// 如果时间为空， 检查
	if expireTime == "" {
		activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, base.Setting.Doubleyear.SantaClausCardId)
		if err != nil {
			return false, err
		}
		if activityRoles.Id == 0 {
			return false, errors.New("activityRoles not fond")
		}
		//保存redis
		key := SANTA_CLAUS_ROLES + userId
		data_map := make(map[string]interface{})
		data_map["itemId"] = activityRoles.ItemId
		data_map["ExpireTime"] = activityRoles.ExpireTime.Format("2006-01-02 15:04:05")
		data_map["RoleId"] = activityRoles.RoleId
		data_map["timestamp"] = activityRoles.ExpireTime.Unix() // 时间戳格式时间
		db.RedisMgr.GetRedisClient().HMSet(key, data_map).Result()
		db.RedisMgr.Expire(key, 3600*24*3) // 过期时间设置3两天
		expireTime = activityRoles.ExpireTime.Format("2006-01-02 15:04:05")
	}
	//判断时间是否已经过期
	currentTime := time.Now().Unix()
	loc, _ := time.LoadLocation("Local") //获取时区
	expTime, _ := time.ParseInLocation("2006-01-02 15:04:05", expireTime, loc)
	exDate := expTime.Unix()
	if currentTime > exDate {
		return false, errors.New("currentTime > expTime")
	}
	return true, nil
}

// 添加过期时间角色过期时间
func (this *SantaClaus) AddRolseExpired(userId string, itmesId int) (int64, error) {
	activityRoles := model.ActivityRoles{
		UserId:     userId,
		ItemId:     itmesId,
		ExpireTime: time.Now().AddDate(0, 0, base.Setting.Doubleyear.SantaClausUnlockEffectiveDay),
	}
	//添加redis
	key := SANTA_CLAUS_ROLES + userId
	data_map := make(map[string]interface{})
	data_map["itemId"] = activityRoles.ItemId
	data_map["ExpireTime"] = activityRoles.ExpireTime.Format("2006-01-02 15:04:05")
	data_map["timestamp"] = activityRoles.ExpireTime.Unix() // 时间戳格式时间
	data_map["RoleId"] = activityRoles.RoleId
	db.RedisMgr.GetRedisClient().HMSet(key, data_map).Result()
	db.RedisMgr.Expire(key, 3600*24*3) // 过期时间设置3两天
	//保存时间
	id, err := db_service.ActivityRolesIns.AddActivityRoles(activityRoles)
	if err != nil {
		logger.Errorf("AddRolseExpired userId=", userId, " itmesId=", itmesId, " err=", err.Error())
		return 0, err
	}
	return id, nil
}

// 移除用户角色
func (this *SantaClaus) RemoveUserRole(userId string) (bool, error) {
	user_info := db.RedisMgr.HGetAll(userId)
	if user_info == nil {
		logger.Errorf("redis Get UserInfo  nil, userId=", userId)
		return false, errors.New("redis err")
	}
	role_id := user_info["role_id"] //  获取角色ID
	roleId, _ := strconv.Atoi(role_id)
	availableRoles := user_info["available_roles"]
	deblockedRoles := user_info["deblocked_roles"]
	//角色code
	cardId := base.Setting.Doubleyear.SantaClausRoleId
	cardIdStr := strconv.Itoa(cardId)
	// 判断当前用户角色是不是圣诞老人
	if roleId == base.Setting.Doubleyear.SantaClausRoleId || strings.Contains(availableRoles, cardIdStr) || strings.Contains(deblockedRoles, cardIdStr) {
		dataMap := map[string]interface{}{}
		if roleId == base.Setting.Doubleyear.SantaClausRoleId {
			// 获取之前角色ID
			activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, base.Setting.Doubleyear.SantaClausCardId)
			if err == nil {
				dataMap["role_id"] = activityRoles.RoleId
				// 更新redis
				db.RedisMgr.GetRedisClient().HSet(userId, "role_id", activityRoles.RoleId)
			}

		}
		if strings.Contains(availableRoles, cardIdStr) {
			availableRolesStr := ""
			availableRolesSlice := strings.Split(availableRoles, "|")
			for _, v := range availableRolesSlice {
				santa_claus_card_id, _ := strconv.Atoi(v)
				if cardId == santa_claus_card_id {
					continue
				}
				availableRolesStr = availableRolesStr + v + "|"
			}
			availableRolesStr = availableRolesStr[:len(availableRolesStr)-1]
			dataMap["available_roles"] = availableRolesStr
			// 更新redis
			db.RedisMgr.GetRedisClient().HSet(userId, "available_roles", availableRolesStr)
		}

		//deblockedRoles
		if strings.Contains(deblockedRoles, cardIdStr) {
			deblockedRolesStr := ""
			deblockedRolesSlice := strings.Split(deblockedRoles, "|")
			for _, v := range deblockedRolesSlice {
				santa_claus_card_id, _ := strconv.Atoi(v)
				if cardId == santa_claus_card_id {
					continue
				}
				deblockedRolesStr = deblockedRolesStr + v + "|"
			}
			deblockedRolesStr = deblockedRolesStr[:len(deblockedRolesStr)-1]
			dataMap["deblocked_roles"] = deblockedRolesStr
			// 更新redis
			db.RedisMgr.GetRedisClient().HSet(userId, "deblocked_roles", deblockedRolesStr)
		}
		if len(dataMap) > 0 {
			//切换到之前的角色
			_, err := db_service.UpdateFields(db_service.UserTable, "user_id", userId, dataMap)
			if err != nil {
				logger.Errorf("RemoveSantaLausRole Update user deblocked_roles err:=", err.Error(), " userId=", userId)
				return false, err
			}
		}

	}

	return true, nil
}

// 移除已解锁的圣诞老人角色卡
func (this *SantaClaus) RemoveSantaLausDeblockedRoles(userId string) (bool, error) {
	// 检查用户是否拥有圣诞老人角色
	user_info := db.RedisMgr.HGetAll(userId)
	if user_info == nil {
		logger.Errorf("redis Get UserInfo  nil, userId=", userId)
		return false, errors.New("redis err")
	}
	// 获取用户已解锁的角色
	deblocked_roles_str := ""
	deblocked_roles := strings.Split(user_info["deblocked_roles"], "|")
	isDeblockedRoles := false
	for _, v := range deblocked_roles {
		roleId, _ := strconv.Atoi(v)
		if base.Setting.Doubleyear.SantaClausRoleId == roleId {
			isDeblockedRoles = true
			break
		}
		deblocked_roles_str = deblocked_roles_str + v + "|"
	}
	if isDeblockedRoles == false {
		return true, nil
	}
	deblocked_roles_str = deblocked_roles_str[:len(deblocked_roles_str)-1] // 去除最后的 |
	// 更新数据库
	dataMap := map[string]interface{}{
		"deblocked_roles": deblocked_roles_str,
	}
	_, err := db_service.UpdateFields(db_service.UserTable, "user_id", userId, dataMap)
	if err != nil {
		logger.Errorf("RemoveSantaLausRole Update user deblocked_roles err:=", err.Error(), " userId=", userId)
		return false, err
	}
	// 更新redis
	db.RedisMgr.GetRedisClient().HSet(userId, "deblocked_roles", deblocked_roles_str)
	return true, nil
}

//移除背包
func (this *SantaClaus) RemoveUserKnapsack(userId string) (bool, error) {
	// 检查用户是否拥有圣诞老人角色
	user_info := db.RedisMgr.HGetAll(userId)
	if user_info == nil {
		logger.Errorf("redis Get UserInfo  nil, userId=", userId)
		return false, errors.New("redis err")
	}
	// 获取用户已解锁的角色
	item_info_str := user_info["item_info"]
	itemId := base.Setting.Doubleyear.SantaClausCardId // 圣诞老人卡片ID
	if strings.Contains(item_info_str, strconv.Itoa(itemId)) {
		// 删除背包的圣诞老人卡
		db_service.UserKnapsackIns.Delete(userId, itemId)
		userKnapsack := make(map[int]map[int]int, 0)
		err := json.Unmarshal([]byte(item_info_str), &userKnapsack)
		if err != nil {
			logger.Errorf("RemoveUserKnapsack, json.Unmarsha err=", err.Error())
			return false, err
		}
		// 重新更新用户背包信息
		userItemMap := make(map[int]map[int]int, 0)
		for key, value := range userKnapsack {
			for k, _ := range value {
				if k != itemId {
					userItemMap[key] = value
				}
			}
		}
		userItemInfo, _ := json.Marshal(userItemMap)
		_, err = db.RedisMgr.GetRedisClient().HSet(userId, "item_info", userItemInfo).Result()
	}
	return true, nil
}

func (this *SantaClaus) GetRoleExpireTime(userId string) string {
	expireTime, _ := db.RedisMgr.GetRedisClient().HGet(SANTA_CLAUS_ROLES+userId, "ExpireTime").Result()
	return expireTime
}
