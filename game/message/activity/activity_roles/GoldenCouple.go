package activity_roles

import (
	"errors"
	"game_server/core/base"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/model"
	"strconv"
	"time"
)

type GoldenCouple struct {
	Roles *ActivityRoles
}

func NewGoldenCouple() *GoldenCouple {
	return &GoldenCouple{
		Roles: NewActivityRoles(),
	}
}

// 判断是否是 金童玉女角色
func (this *GoldenCouple) CheckIsGoldenCoupleRole(roleId int) bool {
	roles := base.Setting.Springfestival.GoldenCoupleRole
	for _, value := range roles {
		if roleId == value {
			return true
		}
	}
	return false
}

// 获取角色
func (this *GoldenCouple) GetRestoreRoleId(roleId int) int {
	roles := base.Setting.ActivityRoles
	for _, value := range roles {
		if roleId == value["role_id"] {
			return value["reset_role_id"].(int)
		}
	}
	return 0
}

// 获取金童玉女的过期天数
func (this *GoldenCouple) GetGoldenCoupleExpireDay(roleId int) int {
	roles := base.Setting.ActivityRoles
	for _, value := range roles {
		if roleId == value["role_id"] {
			return value["usage_days"].(int)
		}
	}
	return 0
}

// 添加 金童玉女过期时间
func (this *GoldenCouple) AddGoldenCoupleRoleExpire(userId string, roleId int) (bool, error) {
	// 检查角色ID
	isGoldenCouple := this.CheckIsGoldenCoupleRole(roleId)
	if isGoldenCouple == false {
		return false, errors.New("not an activity id")
	}
	//检查数据库是否存在金童玉女的数据
	activtiyRoleData, err := db_service.ActivityRolesIns.GetDataByUserId(userId, this.Roles.GetItemIdByRoleId(roleId))
	if err != nil {
		return false, err
	}
	if activtiyRoleData.Id > 0 {
		return true, nil
	}

	//过期时间
	expireTime := time.Now().AddDate(0, 0, this.GetGoldenCoupleExpireDay(roleId))
	var data []model.ActivityRoles
	goldenCoupleRoles := base.Setting.Springfestival.GoldenCoupleRole
	for _, role_id := range goldenCoupleRoles {
		activityRoles := model.ActivityRoles{
			UserId:     userId,
			ItemId:     this.Roles.GetItemIdByRoleId(role_id),
			ExpireTime: expireTime,
			RoleId:     0,
		}
		data = append(data, activityRoles)
		//添加 redis
		key := ACTIVITY_ROLE_EXPIRE + userId + ":roleId:" + strconv.Itoa(role_id)
		data_map := make(map[string]interface{})
		data_map["itemId"] = activityRoles.ItemId
		data_map["ExpireTime"] = expireTime.Format("2006-01-02 15:04:05")
		data_map["timestamp"] = expireTime.Unix() // 时间戳格式时间
		data_map["RoleId"] = 0
		db.RedisMgr.GetRedisClient().HMSet(key, data_map).Result()
		db.RedisMgr.Expire(key, 3600*24*3) // 过期时间设置3两天
	}
	//保存时间
	return db_service.ActivityRolesIns.BatchAdd(data)
}

func (this *GoldenCouple) GetRoleExpireTime(userId string, roleId int) string {
	key := ACTIVITY_ROLE_EXPIRE + userId + ":roleId:" + strconv.Itoa(roleId)
	expireTime, _ := db.RedisMgr.GetRedisClient().HGet(key, "ExpireTime").Result()
	return expireTime
}

// 记录用户就的ID
func (this *GoldenCouple) SaveUserOldRole(userId string, roleId int) (int64, error) {
	var itmesId []int
	roles := base.Setting.Springfestival.GoldenCoupleRole
	for _, value := range roles {
		itmesId = append(itmesId, this.Roles.GetItemIdByRoleId(value))
	}
	return db_service.ActivityRolesIns.ModifyGoldenCoupleOldUserRole(userId, roleId, itmesId)
}

// 获取解锁
func (this *GoldenCouple) GetUnloc() string {
	//更新redis
	RoleIds := base.Setting.Springfestival.GoldenCoupleRole
	RoleIdStr := ""
	for _, id := range RoleIds {
		value := strconv.Itoa(id)
		RoleIdStr += "|" + value
	}
	return RoleIdStr
}
