package db_service

import (
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/model"
	"time"
)

type ActivityInvite struct {
}

const (
	USER_INVITATION string = "t_activity_invitation"
)

func (this *ActivityInvite) AddUserInvit(userId string, inviterId string) (int64, error) {
	session := db.Mysql.Table(USER_INVITATION)
	userInvit := model.ActivityInvitation{}
	userInvit.UserId = userId
	userInvit.InviterId = inviterId
	userInvit.IsCreateRole = 0
	userInvit.UpdateTime = time.Now()
	_, err := session.Insert(&userInvit)
	if err != nil {
		return 0, err
	}
	return userInvit.Id, nil
}

//获取用户邀请列表
func (this *ActivityInvite) GetUserInvitList(userId string) ([]model.ActivityUserInvitList, error) {
	list := make([]model.ActivityUserInvitList, 0)
	session := db.Mysql.Table(USER_INVITATION).Join("INNER", UserTable, "t_user.user_id = t_activity_invitation.user_id")
	session.Cols("t_user.user_id,t_user.nick_name,t_activity_invitation.create_time,t_activity_invitation.update_time")
	// 排行榜开始时间
	session.Where("t_activity_invitation.create_time >= ?", base.Setting.Springfestival.RankingListStartDate+" 00:00:00")
	session.Where("t_activity_invitation.is_create_role >= ?", 1)
	err := session.Where("t_activity_invitation.inviter_id = ?", userId).Find(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

//获取用户邀请列表
func (this *ActivityInvite) GetUserDayInvitList(userId string, startTime, endTime string) ([]model.ActivityUserInvitList, error) {
	list := make([]model.ActivityUserInvitList, 0)
	session := db.Mysql.Table(USER_INVITATION).Join("INNER", UserTable, "t_user.user_id = t_activity_invitation.user_id")
	session.Cols("t_user.user_id,t_user.nick_name,t_activity_invitation.create_time,t_activity_invitation.update_time")
	session.Where("t_activity_invitation.inviter_id = ?", userId)
	session.Where("t_activity_invitation.is_create_role >= ?", 1)
	session.Where("t_activity_invitation.update_time >= ?", startTime)
	session.Where("t_activity_invitation.update_time <= ?", endTime)
	err := session.Find(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

//获取用户的邀请信息
func (this *ActivityInvite) GetInvitInfoByUserId(userId string) (model.ActivityInvitation, error) {
	var invitInfo model.ActivityInvitation
	session := db.Mysql.Table(USER_INVITATION).Where("user_id = ?", userId)
	_, err := session.Get(&invitInfo)
	if err != nil {
		logger.Errorf("GetInvitInfoByUserId: %s, err:%v \n", userId, err)
		return invitInfo, err
	}
	return invitInfo, nil
}

//更新用户是否创建角色
func (this *ActivityInvite) UpdateIsCreateRole(userId string, isCreateRole int) (int64, error) {
	session := db.Mysql.Table(USER_INVITATION).Where("user_id = ?", userId)
	session.Where("is_create_role = ?", 0)
	data := map[string]interface{}{"is_create_role": isCreateRole, "update_time": time.Now().Format("2006-01-02 15:04:05")}
	return session.Update(data)
}

/**
获取用户邀请数
*/
func (this *ActivityInvite) GetUserDayInvitCount(userId string, startTime, endTime string) (int64, error) {
	session := db.Mysql.Table(USER_INVITATION)
	session.Where("inviter_id = ?", userId)
	session.And("is_create_role = 1")
	session.And("create_time >= ?", startTime)
	session.And("create_time <= ?", endTime)
	num, err := session.Count()
	if err != nil {
		return num, err
	}
	return num, nil
}

func (this *ActivityInvite) GetUserOldInvitCount(userId string, endTime string, activity_beginTime string) (int64, error) {
	session := db.Mysql.Table(USER_INVITATION)
	session.Where("inviter_id = ?", userId)
	session.And("is_create_role = 1")
	session.And("create_time <= ?", endTime)
	session.And("create_time >= ?", activity_beginTime)
	num, err := session.Count()
	if err != nil {
		return num, err
	}
	return num, nil
}

/**
获取总用户邀请数
*/
func (this *ActivityInvite) GetUserInvitCount(userId string) (int64, error) {
	session := db.Mysql.Table(USER_INVITATION)
	session.Where("inviter_id = ?", userId)
	num, err := session.Count()
	if err != nil {
		return num, err
	}
	return num, nil
}
