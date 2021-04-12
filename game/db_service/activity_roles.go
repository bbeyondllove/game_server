package db_service

import (
	"game_server/core/base"
	"game_server/db"
	"game_server/game/model"
)

type ActivityRoles struct {
}

const (
	ACTIVITY_ROLES string = "t_activity_roles"
)

// 添加活动角色
func (this *ActivityRoles) AddActivityRoles(date_map model.ActivityRoles) (int64, error) {
	_, err := db.Mysql.Insert(&date_map)
	if err != nil {
		return 0, err
	}
	return date_map.Id, nil
}

// 批量添加
func (this *ActivityRoles) BatchAdd(data []model.ActivityRoles) (bool, error) {
	_, err := db.Mysql.Insert(&data)
	if err != nil {
		return false, err
	}
	return true, nil
}

//通过用户id， 获取活动角色
func (this *ActivityRoles) GetDataByUserId(userId string, itemId int) (model.ActivityRoles, error) {
	activityRoles := model.ActivityRoles{}
	_, err := db.Mysql.Table(ACTIVITY_ROLES).Where("user_id = ?", userId).Where("item_id = ?", itemId).OrderBy("id desc").Limit(1).Get(&activityRoles)
	return activityRoles, err
}

func (this *ActivityRoles) ModifyUserRoleId(userId string, roleId int) (int64, error) {
	data := map[string]interface{}{
		"role_id": roleId,
	}
	return db.Mysql.Table(ACTIVITY_ROLES).Where("user_id = ?", userId).Where("item_id = ?", base.Setting.Doubleyear.SantaClausCardId).Update(data)
}

// 更新
func (this *ActivityRoles) ModifyGoldenCoupleOldUserRole(userId string, roleId int, itemIs []int) (int64, error) {
	data := map[string]interface{}{
		"role_id": roleId,
	}
	return db.Mysql.Table(ACTIVITY_ROLES).Where("user_id = ?", userId).In("item_id", itemIs).Update(data)
}
