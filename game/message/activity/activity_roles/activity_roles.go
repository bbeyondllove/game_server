package activity_roles

import (
	"errors"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/message/activity/santa_claus"
	"game_server/game/model"
	"strconv"
	"time"
)

const ACTIVITY_ROLE_EXPIRE = "activity_roles_expire_time:" // 角色过期时间

//活动角色
type ActivityRoles struct {
}

func NewActivityRoles() *ActivityRoles {
	return &ActivityRoles{}
}

//检查是否是活动角色
func (this *ActivityRoles) CheckIsActivityRole(roleId int) bool {
	// 兼容圣诞老人角色
	if roleId == base.Setting.Doubleyear.SantaClausRoleId {
		return true
	}
	// 新的活动
	rolesConf := base.Setting.ActivityRoles
	for _, conf := range rolesConf {
		if roleId == conf["role_id"].(int) {
			return true
		}
	}
	return false
}

//检查是否是活动角色卡
func (this *ActivityRoles) CheckIsActivityRoleCard(cardId int) bool {
	// 兼容圣诞老人角色
	if cardId == base.Setting.Doubleyear.SantaClausCardId {
		return true
	}
	// 新的活动
	rolesConf := base.Setting.ActivityRoles
	for _, conf := range rolesConf {
		if cardId == conf["role_card_id"].(int) {
			return true
		}
	}
	return false
}

//检查活动角色是否可以解锁  true 可以使用   false 不可以使用
func (this *ActivityRoles) CheckRoleCardIsUnLock(cardId int) bool {
	// 兼容圣诞老人
	if cardId == base.Setting.Doubleyear.SantaClausCardId {
		isExpire, _ := santa_claus.SantaClausLogic.CheckEndExpire()
		if isExpire == false {
			return false
		}
	}
	// 新的活动
	rolesConf := base.Setting.ActivityRoles
	for _, conf := range rolesConf {
		if cardId == conf["role_card_id"].(int) {
			//设置时区
			loc, _ := time.LoadLocation("Local")
			timeNow := time.Now().Unix()
			roleStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", conf["role_start_time"].(string), loc) // 开始时间，
			roleEndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", conf["role_end_time"].(string), loc)     //  结束时间
			if timeNow >= roleStartTime.Unix() && timeNow <= roleEndTime.Unix() {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// 检查角色是否过期
// false 已过期  true 可以选择
func (this *ActivityRoles) CheckRolesIsExpire(userId string, roleId int) (bool, error) {
	//兼容之前老的圣诞老人
	if roleId == base.Setting.Doubleyear.SantaClausRoleId {
		isExpire, err := santa_claus.SantaClausLogic.CheckUserRoleExpire(userId) // 检查角色是否过期
		if isExpire == false {
			return isExpire, err
		}
	}
	key := ACTIVITY_ROLE_EXPIRE + userId + ":roleId:" + strconv.Itoa(roleId)
	expireTime, _ := db.RedisMgr.GetRedisClient().HGet(key, "ExpireTime").Result()
	// 如果时间为空， 检查
	if expireTime == "" {
		activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, this.GetItemIdByRoleId(roleId))
		if err != nil {
			return false, err
		}
		if activityRoles.Id == 0 {
			return false, errors.New("activityRoles not fond")
		}
		//保存redis
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
		return false, nil
	}
	return true, nil
}

// 根据角色ID获取ItemId
func (this *ActivityRoles) GetItemIdByRoleId(roleId int) int {
	rolesConf := base.Setting.ActivityRoles
	for _, conf := range rolesConf {
		if conf["role_id"].(int) == roleId {
			return conf["role_card_id"].(int)
		}
	}
	return 0
}

// 添加过期时间角色过期时间
func (this *ActivityRoles) AddRolseExpired(userId string, roleId int, itmesId int) (int64, error) {
	activityRoles := model.ActivityRoles{
		UserId:     userId,
		ItemId:     itmesId,
		ExpireTime: time.Now().AddDate(0, 0, base.Setting.Doubleyear.SantaClausUnlockEffectiveDay),
	}
	//添加redis
	key := ACTIVITY_ROLE_EXPIRE + userId + ":roleId:" + strconv.Itoa(roleId)
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
