package message

import (
	"errors"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/message/activity/activity_roles"
	"strconv"
	"strings"
)

// 检测用户
func CheckColdenCouple(userId string, roleId int, availableRoles string, deblockedRoles string) (int, string, string) {
	userRoleId := roleId
	userAvailableRoles := availableRoles
	userDeblockedRoles := deblockedRoles
	// 检查是否是金童玉女角色
	activityRole := activity_roles.NewGoldenCouple()
	isColdenCoupleRole := activityRole.CheckIsGoldenCoupleRole(userRoleId)

	if isColdenCoupleRole || checkDeblockedRoles(deblockedRoles) {
		//判断是否已经过期
		// 检查角色是否过期是否可用
		goldenCoupleRole := roleId
		goldenCouple := base.Setting.Springfestival.GoldenCoupleRole
		if isColdenCoupleRole == false {
			goldenCoupleRole = goldenCouple[0]
		}
		isExpire, _ := activityRole.Roles.CheckRolesIsExpire(userId, goldenCoupleRole)
		if isExpire == false {
			activityRoles, err := db_service.ActivityRolesIns.GetDataByUserId(userId, activityRole.Roles.GetItemIdByRoleId(goldenCoupleRole))
			if err == nil {
				userRoleId = activityRoles.RoleId
			} else {
				//配置的复原ID
				userRoleId = activityRole.GetRestoreRoleId(roleId)
			}

			// 检查解锁列表是否包含金童玉女
			if checkDeblockedRoles(availableRoles) {
				availableRolesStr := ""
				availableRolesSlice := strings.Split(availableRoles, "|")
				for _, v := range availableRolesSlice {
					golden_couple_role_id, _ := strconv.Atoi(v)
					if golden_couple_role_id == goldenCouple[0] || golden_couple_role_id == goldenCouple[1] {
						continue
					}
					availableRolesStr = availableRolesStr + v + "|"
				}
				userAvailableRoles = availableRolesStr[:len(availableRolesStr)-1]
			}

			//deblockedRoles 检查列表是否包含金童玉女
			if checkDeblockedRoles(deblockedRoles) {
				deblockedRolesStr := ""
				deblockedRolesSlice := strings.Split(deblockedRoles, "|")
				for _, v := range deblockedRolesSlice {
					golden_couple_role_id, _ := strconv.Atoi(v)
					if golden_couple_role_id == goldenCouple[0] || golden_couple_role_id == goldenCouple[1] {
						continue
					}
					deblockedRolesStr = deblockedRolesStr + v + "|"
				}
				userDeblockedRoles = deblockedRolesStr[:len(deblockedRolesStr)-1]
			}

			email := make(map[string]interface{})
			email["userId"] = userId
			email["emailType"] = 1
			email["emailTitle"] = "金童玉女已过期"
			email["emailContent"] = "您的金童玉女职业解锁卡（3天）已过期。"
			db_service.EmailLogicIns.AddEmail(email)
		}
	}
	return userRoleId, userAvailableRoles, userDeblockedRoles
}

// 检查用户以解锁列表是否包含活动角色
func checkDeblockedRoles(roleList string) bool {
	goldenCouple := base.Setting.Springfestival.GoldenCoupleRole
	for _, v := range goldenCouple {
		str := strconv.Itoa(v)
		if strings.Contains(roleList, str) == true {
			return true
		}
	}
	return false
}

// 移除过期金童玉女
func RemoveExpireColdenCoupl(userId string) (bool, error) {

	user_info := db.RedisMgr.HGetAll(userId)
	if user_info == nil {
		logger.Errorf("redis Get UserInfo  nil, userId=[%v]", userId)
		return false, errors.New("redis err")
	}
	role_id := user_info["role_id"] //  获取角色ID
	roleId, _ := strconv.Atoi(role_id)
	coldenCouplRole := roleId
	availableRoles := user_info["available_roles"]
	deblockedRoles := user_info["deblocked_roles"]
	activityRole := activity_roles.NewGoldenCouple()
	if activityRole.CheckIsGoldenCoupleRole(roleId) == false {
		return false, nil
	}
	//检查是否过期
	isExpire, err := activityRole.Roles.CheckRolesIsExpire(userId, roleId)
	if err != nil {
		return false, err
	}

	if isExpire == true {
		go ColdenCoupleRedEnvelope(userId)
		return false, nil
	}
	//
	roleId, availableRoles, deblockedRoles = CheckColdenCouple(userId, roleId, availableRoles, deblockedRoles)
	// 判断是否是活动ID
	if activityRole.CheckIsGoldenCoupleRole(roleId) == true {
		return false, nil
	}
	//
	dataMap := map[string]interface{}{}
	dataMap["role_id"] = roleId
	dataMap["available_roles"] = availableRoles
	dataMap["deblocked_roles"] = deblockedRoles
	_, err = db_service.UpdateFields(db_service.UserTable, "user_id", userId, dataMap)
	if err != nil {
		return false, err
	}
	db.RedisMgr.GetRedisClient().HMSet(userId, dataMap)
	// 通知用户已过期
	go NotifyUsersRoleIsExpire(userId, activityRole.Roles.GetItemIdByRoleId(coldenCouplRole))
	return true, nil
}
