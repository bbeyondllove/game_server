package message

import (
	"errors"
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/message/activity/santa_claus"
	"game_server/game/proto"
	"strconv"
	"strings"
)

func CheckCurrentRole(userId string) (bool, error) {
	//获取当前用户的角色id
	role_id, err := db.RedisMgr.GetRedisClient().HGet(userId, "role_id").Result()
	if err != nil {
		logger.Errorf("get user role_id get redis err:=", err.Error(), "userId = ", userId)
		return false, err
	}
	roleId, _ := strconv.Atoi(role_id)
	// 判断当前用户角色是不是圣诞老人
	if roleId != base.Setting.Doubleyear.SantaClausRoleId {
		return false, errors.New("It's not Santa Claus")
	}
	// 检查角色是否过期是否可用
	isExpire, err := santa_claus.SantaClausLogic.CheckUserRoleExpire(userId)
	if isExpire == false {
		// 移除用户角色
		santa_claus.SantaClausLogic.RemoveUserRole(userId)
		// 通知用户已过期
		go NotifyUsersRoleIsExpire(userId, base.Setting.Doubleyear.SantaClausCardId)
		email := make(map[string]interface{})
		email["userId"] = userId
		email["emailType"] = 1
		email["emailTitle"] = "圣诞老人已过期"
		email["emailContent"] = "您的圣诞老人职业解锁卡（3天）已过期。"
		db_service.EmailLogicIns.AddEmail(email)
		return false, err
	}
	return isExpire, err
}

//通知前端更
func NotifyUsersRoleIsExpire(userId string, itemId int) {
	// 获取用户之前的角色ID
	//key := santa_claus.SANTA_CLAUS_ROLES + userId
	//roleId_str, err := db.RedisMgr.GetRedisClient().HGet(key, "RoleId").Result()
	//if err != nil {
	//	logger.Errorf("NotifyUsersRoleIsExpire get redis err:=", err.Error(), "userId=", userId)
	//	return
	//}
	//roleId, _ := strconv.Atoi(roleId_str)
	//if roleId == 0 {
	//	logger.Errorf("NotifyUsersRoleIsExpire roleId == 0  userId:=", userId)
	//	return
	//}
	activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, itemId)
	if err != nil {
		logger.Errorf("NotifyUsersRoleIsExpire push faild userid:[%v], err= %v", userId, err)
	}
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SANTA_LAUS_ROLE_EXPIRE_RSP)
	pushMessage := &proto.S2CSantaClausRoleExprie{}
	pushMessage.Code = errcode.MSG_SUCCESS
	pushMessage.UserId = userId
	pushMessage.RoleId = activityRoles.RoleId
	fmt.Println(pushMessage)
	logger.Debugf("NotifyUsersRoleIsExpire end")
	rsp.WriteData(pushMessage)
	Sched.SendToAllUser(rsp)
}

// 检测用户的圣诞老人卡
func CheckSantaClaus(userId string, roleId int, availableRoles string, deblockedRoles string) (int, string, string) {
	userRoleId := roleId
	userAvailableRoles := availableRoles
	userDeblockedRoles := deblockedRoles
	//圣诞老人
	settingRoleId := base.Setting.Doubleyear.SantaClausRoleId //
	roleIdStr := strconv.Itoa(settingRoleId)
	if roleId == settingRoleId || strings.Contains(deblockedRoles, roleIdStr) {
		//判断是否已经过期
		// 检查角色是否过期是否可用
		isExpire, _ := santa_claus.SantaClausLogic.CheckUserRoleExpire(userId)
		if isExpire == false {
			if roleId == base.Setting.Doubleyear.SantaClausRoleId {
				//获取用户之前的RoleId
				activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, base.Setting.Doubleyear.SantaClausCardId)
				if err == nil {
					userRoleId = activityRoles.RoleId
				}
			}

			// availableRoles
			if strings.Contains(availableRoles, roleIdStr) {
				availableRolesStr := ""
				availableRolesSlice := strings.Split(availableRoles, "|")
				for _, v := range availableRolesSlice {
					santa_claus_card_id, _ := strconv.Atoi(v)
					if settingRoleId == santa_claus_card_id {
						continue
					}
					availableRolesStr = availableRolesStr + v + "|"
				}
				userAvailableRoles = availableRolesStr[:len(availableRolesStr)-1]
			}

			//deblockedRoles
			if strings.Contains(deblockedRoles, roleIdStr) {
				deblockedRolesStr := ""
				deblockedRolesSlice := strings.Split(deblockedRoles, "|")
				for _, v := range deblockedRolesSlice {
					santa_claus_card_id, _ := strconv.Atoi(v)
					if settingRoleId == santa_claus_card_id {
						continue
					}
					deblockedRolesStr = deblockedRolesStr + v + "|"
				}
				userDeblockedRoles = deblockedRolesStr[:len(deblockedRolesStr)-1]
			}

			// 更新数据库
			//dataMap := map[string]interface{}{
			//	"role_id":         userRoleId,
			//	"available_roles": userAvailableRoles,
			//	"deblocked_roles": userDeblockedRoles,
			//}
			//db_service.UpdateFields(db_service.UserTable, "user_id", userId, dataMap)
			// 添加邮件
			email := make(map[string]interface{})
			email["userId"] = userId
			email["emailType"] = 1
			email["emailTitle"] = "圣诞老人已过期"
			email["emailContent"] = "您的圣诞老人职业解锁卡（3天）已过期。"
			db_service.EmailLogicIns.AddEmail(email)
		}
	}
	return userRoleId, userAvailableRoles, userDeblockedRoles
}

//检测背包中的圣诞老人，过期之后删除
func CheckUserKnapsacSantaClaus(userId string, itemData map[int]map[int]int) map[int]map[int]int {
	userItemMap := make(map[int]map[int]int, 0)
	for key, value := range itemData {
		for k, _ := range value {
			if k == base.Setting.Doubleyear.SantaClausCardId {
				//检测是否已过期
				expire, _ := santa_claus.SantaClausLogic.CheckEndExpire()
				if expire == false {
					//删除背包
					db_service.UserKnapsackIns.Delete(userId, k)
					break
				}
			}
			userItemMap[key] = value

		}
	}
	return userItemMap
}
