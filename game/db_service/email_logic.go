package db_service

import (
	"encoding/json"
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/errcode"
	"game_server/game/model"
	"strconv"
	"time"
)

var EmailLogicIns = NewEmailLogic()

type EmailLogic struct {
}

const EMAIL_PRIZE_KEY = "email:prize:emailId:"

func (this *EmailLogic) GetEmailPrize(emailId int64) ([]model.EmailPrize, error) {
	key := EMAIL_PRIZE_KEY + strconv.FormatInt(emailId, 10)
	result := db.RedisMgr.Get(key)
	var pirzeList = make([]model.EmailPrize, 0)
	if len(result) > 0 {
		err := json.Unmarshal([]byte(result), &pirzeList)
		if err != nil {
			logger.Errorf("GetEmailPrize json.Unmarshal error, err=", err.Error())
			return pirzeList, err
		}
		return pirzeList, nil
	}
	//从数据库获取
	pirzeList, err := EmailPrizeIns.GetEmailPrize(emailId)
	if err != nil {
		logger.Errorf("GetEmailPrize db error, err=", err.Error())
		return pirzeList, err
	}
	json, err := json.Marshal(pirzeList)
	if err != nil {
		logger.Errorf("GetEmailPrize json.Marshal, err=", err.Error())
		return pirzeList, err
	}
	//缓存一个小时
	db.RedisMgr.GetRedisClient().Set(key, string(json), time.Second*3600)
	return pirzeList, nil
}

/**
添加邮件
*/
func (this *EmailLogic) AddEmail(emailInfo map[string]interface{}) (int64, bool) {
	var email model.Email
	email.EmailType = emailInfo["emailType"].(int)
	email.UserId = emailInfo["userId"].(string)
	email.EmailTitle = emailInfo["emailTitle"].(string)
	email.EmailContent = emailInfo["emailContent"].(string)
	email.Display = 1 // 邮件显示
	// 普通邮件
	if email.EmailType == 1 {
		email, err := EmailIns.AddEmail(email.UserId, email.EmailType, email.EmailTitle, email.EmailContent)
		if err != nil {
			logger.Errorf("AddEmail db  err=", err)
			return 0, false
		}
		return email.Id, true
	} else if email.EmailType == 2 {
		emailPirze := make([]model.EmailPrize, 0)
		prizeList := emailInfo["prizeList"]
		for _, v := range prizeList.([]map[string]interface{}) {
			temp := model.EmailPrize{}
			temp.PrizeName = v["prizeName"].(string)
			temp.PrizeImg = v["prizeImg"].(string)
			temp.PrizeType = v["prizeType"].(int)
			temp.PrizeId = v["prizeId"].(int)
			temp.PrizeNum = v["prizeNum"].(int)
			if _, ok := v["extend"].(string); ok {
				temp.Extend = v["extend"].(string)
			}
			emailPirze = append(emailPirze, temp)
		}
		id, err := EmailIns.AddEmailPrize(email, emailPirze)
		if err != nil {
			logger.Errorf("AddEmail db  err=", err)
			return 0, false
		}
		return id, true
	} else {
		logger.Errorf("AddEmail EmailType err", email.EmailType)
		return 0, false
	}
}

//领取奖励
func (this *EmailLogic) ReceiveRewards(userId string, emailId int) (int32, model.Email, bool) {
	var email model.Email
	var err error
	email, err = EmailIns.GetEmail(userId, int64(emailId))
	if err != nil { // 领取奖励失败
		logger.Errorf("Email ReceiveRewards  err=", err)
		return errcode.ERROR_EMAIL_FAILED_TO_CLAIM_REWARD, email, false
	}
	if email.Id == 0 { // 邮件不存在
		return errcode.ERROR_EMAIL_NOT_FOUND, email, false
	}
	if email.EmailType != 2 { // 领取奖励失败
		logger.Errorf("Email ReceiveRewards  email.EmailType !=2 :", email.EmailType)
		return errcode.ERROR_EMAIL_FAILED_TO_CLAIM_REWARD, email, false
	}
	if email.IsRead == 1 { //奖励已领取
		return errcode.ERROR_EMAIL_REWARD_RECEIVED, email, false
	}
	// 获取奖品
	prizeList, err := EmailPrizeIns.GetEmailPrize(email.Id)
	if err != nil {
		logger.Errorf("Email ReceiveRewards get przieId  err=", err)
		return errcode.ERROR_EMAIL_FAILED_TO_CLAIM_REWARD, email, false
	}
	if len(prizeList) == 0 {
		logger.Errorf("Email ReceiveRewards Email prize is 0")
		return errcode.ERROR_EMAIL_FAILED_TO_CLAIM_REWARD, email, false
	}

	// 检查是否是活动道具
	if this.checkIsActivityRoleCard(prizeList[0].PrizeId) {
		//检查是否过期
		if this.checkRoleCardIsExpire(prizeList[0].PrizeId) == false {
			return errcode.ERROR_PRIZE_EXPIRED, email, false
		}
	}

	//  发放奖励
	if isSeceive := this.emailPrize(userId, emailId, prizeList); isSeceive != true {
		return errcode.ERROR_EMAIL_FAILED_TO_CLAIM_REWARD, email, false
	}
	//清除缓存
	db.RedisMgr.GetRedisClient().Del(EMAIL_PRIZE_KEY + strconv.Itoa(emailId)).Result()
	return errcode.MSG_SUCCESS, email, true
}

func (this *EmailLogic) emailPrize(userId string, emailId int, prizeList []model.EmailPrize) bool {
	session := db.Mysql
	for _, v := range prizeList {
		// 领取奖励
		number, err := session.Table(EMAIL_PRIZE_TABLE).Where("id = ?", v.Id).Where("is_receive = 0").Update(map[string]int{"is_receive": 1})
		if err != nil {
			logger.Errorf("emailPrize update err=", err.Error())
			return false
		}
		// 更新条数等于 0, 领取失败
		if number == 0 {
			logger.Errorf("emailPrize update number == 0 :", number)
			return false
		}
		_, err = EmailPrizeAction[v.PrizeType].Handler(userId, v)
		fmt.Println()
		// 领取奖励失败
		if err != nil {
			session.Table(EMAIL_PRIZE_TABLE).Where("id = ?", v.Id).Update(map[string]int{"is_receive": 2})
		}
	}
	//更新邮件已读
	session.Table(EMAIL_TABLE).Where("user_id = ? ", userId).Where("id = ?", emailId).Update(map[string]int{"is_read": 1})
	return true
}

/**
重新加载用户背包
*/
func (this *EmailLogic) ReloadUserKnapsack(userId string, emailPrize model.EmailPrize) (bool, error) {
	//查询奖品
	item, _ := ItemIns.GetDataById(emailPrize.PrizeId)
	// 新增背包记录
	userdata := model.UserKnapsack{
		UserId:   userId,
		ItemType: item.ItemType,
		ItemId:   item.Id,
		ItemNum:  emailPrize.PrizeNum,
	}
	_, err := UserKnapsackIns.Add(&userdata)
	if err != nil {
		logger.Errorf("emailPrize db_service.UserKnapsackIns.Add err=", err.Error())
		return false, err
	}
	userKnapsack, err := UserKnapsackIns.GetDataByUid(userId)
	if err != nil {
		logger.Errorf("reloadUserKnapsack failed, err=", err.Error())
		return false, err
	}
	userItemMap := make(map[int]map[int]int, 0)
	for _, v := range userKnapsack {
		if _, ok := userItemMap[v.ItemType]; !ok {
			userItemMap[v.ItemType] = make(map[int]int)
		}
		userItemMap[v.ItemType][v.ItemId] = v.ItemNum
	}
	userItemInfo, _ := json.Marshal(userItemMap)
	_, err = db.RedisMgr.GetRedisClient().HSet(userId, "item_info", userItemInfo).Result()
	AvailableRoles, _ := db.RedisMgr.GetRedisClient().HGet(userId, "available_roles").Result()
	if AvailableRoles != "" {
		//更新redis
		RoleIds := base.Setting.Springfestival.GoldenCoupleRole
		RoleIdStr := ""
		for _, id := range RoleIds {
			value := strconv.Itoa(id)
			RoleIdStr += "|" + value
		}

		AvailableRoles = AvailableRoles + RoleIdStr
		_, err = db.RedisMgr.GetRedisClient().HSet(userId, "available_roles", AvailableRoles).Result()
		// 更新db
		dataMap := make(map[string]interface{})
		dataMap["available_roles"] = AvailableRoles
		UpdateFields(UserTable, "user_id", userId, dataMap)

	}
	// 获取
	return true, nil
}

//双旦活动，领取奖励
func (this *EmailLogic) ActivityReward(userId string, emailPrize model.EmailPrize) (bool, error) {
	extend := emailPrize.Extend //扩展消息
	dataMap := make(map[string]int, 0)
	err := json.Unmarshal([]byte(extend), &dataMap)
	if err != nil {
		return false, err
	}
	code, _ := NewCdt().UpdateUserCdt(userId, float32(emailPrize.PrizeNum), dataMap["eventType"])
	logger.Debugf("ActivityReward, UpdateUserCdt code=%s, emailId =%s", code, emailPrize.EmailId)
	return true, nil
}

//检查是否是活动角色卡
func (this *EmailLogic) checkIsActivityRoleCard(cardId int) bool {
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

//检查活动角色卡是否过期
func (this *EmailLogic) checkRoleCardIsExpire(cardId int) bool {
	// 兼容圣诞老人
	if cardId == base.Setting.Doubleyear.SantaClausCardId {
		//检测是否已过期
		currentTime := time.Now().Unix()                                                                  // 当前时间
		expireTime, _ := time.Parse("2006-01-02 15:04:05", base.Setting.Doubleyear.SantaClausExpriteTime) //  圣诞老人过期司机
		if currentTime > expireTime.Unix() {
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

// NewEmailLogic 实例化EmailLogic结构体.
func NewEmailLogic() *EmailLogic {
	return &EmailLogic{}
}
