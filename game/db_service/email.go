package db_service

import (
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/model"
	"game_server/game/proto"
	"time"
)

type EmailInfo struct{}

const (
	EMAIL_TABLE = "t_email"
)

func (this *EmailInfo) GetEmail(userId string, email_id int64) (model.Email, error) {
	session := db.Mysql.Table(EMAIL_TABLE)
	session.Where("id = ? ", email_id)
	session.Where("user_id = ? and display = 1 ", userId)
	email := model.Email{}
	has, err := session.Get(&email)
	if err != nil || has == false {
		return email, err
	}
	return email, nil
}

// 获取email 数量
func (this *EmailInfo) GetEmailCount(userId string, whereData map[string]interface{}) (int64, error) {
	session := db.Mysql.Table(EMAIL_TABLE)
	session.Where("user_id = ? and display = 1 ", userId)
	session.Where("expire_time >= ?", time.Now().Format("2006-01-02 15:04:05"))
	if len(whereData) > 0 {
		for key, value := range whereData {
			fieldName := fmt.Sprintf(" %s = ?", key)
			session.Where(fieldName, value)
		}
	}
	return session.Count()
}

//获取用户已读未读数量
func (this *EmailInfo) GetReadAndUnreadNum(userId string) (map[string]int64, error) {
	var unreadNum, readNum int64
	ret := make(map[string]int64, 0)
	var err error
	// 未读的邮件
	unreadNum, err = this.GetEmailCount(userId, map[string]interface{}{"is_read": 0})

	if err != nil {
		return nil, err
	}
	// 已读的邮件
	readNum, err = this.GetEmailCount(userId, map[string]interface{}{"is_read": 1})
	if err != nil {
		return nil, err
	}
	ret["unread"] = unreadNum
	ret["read"] = readNum
	return ret, nil
}

//获取用户邮件
func (this *EmailInfo) GetDataByUserId(userid string, isRead int) ([]proto.EmailItem, error) {
	emailList := make([]model.Email, 0)
	session := db.Mysql.Table(EMAIL_TABLE)
	session.Where("expire_time >= ?", time.Now().Format("2006-01-02 15:04:05"))
	session.Where("user_id = ? and display = ?", userid, 1).OrderBy("id desc").Limit(200)
	if isRead != -1 {
		session.Where("is_read = ?", isRead)
	}
	err := session.Find(&emailList)
	if err != nil {
		return nil, err
	}

	readItem := make([]proto.EmailItem, 0)
	for _, email := range emailList {
		temp := proto.EmailItem{
			Id:           email.Id,
			UserId:       email.UserId,
			EmailType:    email.EmailType,
			EmailTitle:   email.EmailTitle,
			EmailContent: email.EmailContent,
			IsRead:       email.IsRead,
			CreateTime:   email.CreateTime.Format("2006/01/02 15:04:05"),
			ExpireTime:   email.ExpireTime.Format("2006/01/02 15:04:05"),
		}
		readItem = append(readItem, temp)
	}
	return readItem, nil

}

// 删除邮件
func (this *EmailInfo) DelEmailInIds(userid string, ids []int) (int64, error) {
	dbSession := db.Mysql.Table(EMAIL_TABLE)
	dbSession.Where(" user_id = ? ", userid)
	dbSession.In("id", ids)
	var email model.Email
	email.Display = 0
	email.UpdateTime = time.Now()
	return dbSession.Cols("display").Update(email)
}

// 设置邮件为已读,普通邮件
func (this *EmailInfo) SetEmailIsRead(userid string, ids []int) (int64, error) {
	dbSession := db.Mysql.Table(EMAIL_TABLE)
	dbSession.Where("user_id = ? ", userid).Where("email_type = ?", 1)
	dbSession.In("id", ids)
	var email model.Email
	email.IsRead = 1
	email.UpdateTime = time.Now()
	return dbSession.Update(email)
}

//添加邮件
func (this *EmailInfo) AddEmail(userId string, emailType int, emailTilte string, emailText string) (model.Email, error) {
	var email model.Email
	email.Display = 1
	email.UserId = userId
	email.EmailType = emailType
	email.EmailTitle = emailTilte
	email.EmailContent = emailText
	email.IsRead = 0
	email.CreateTime = time.Now()
	email.UpdateTime = time.Now()
	// 设置邮箱过期时间
	email.ExpireTime = time.Now().AddDate(0, 0, base.Setting.Email.EmailExpireDay)
	dbSession := db.Mysql.Table(EMAIL_TABLE)
	_, err := dbSession.Insert(&email)
	if err != nil {
		return email, err
	}
	return email, nil
}

/**
添加带奖励的邮件
*/
func (this *EmailInfo) AddEmailPrize(email model.Email, emailPirze []model.EmailPrize) (int64, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		logger.Errorf("AddEmailPrize session.Begin err", err)
		return 0, err
	}
	email.CreateTime = time.Now()
	email.UpdateTime = time.Now()
	// 设置邮箱过期时间
	email.ExpireTime = time.Now().AddDate(0, 0, base.Setting.Email.EmailExpireDay)
	_, err = session.Table(EMAIL_TABLE).Insert(&email)
	if err != nil {
		logger.Errorf("AddEmailPrize session.Insert err", err)
		session.Rollback()
		return 0, err
	}
	for k, _ := range emailPirze {
		emailPirze[k].EmailId = int(email.Id)
	}
	_, err = session.Table(EMAIL_PRIZE_TABLE).Insert(emailPirze)
	if err != nil {
		logger.Errorf("AddEmailPrize EmailPrize session.Insert err", err)
		session.Rollback()
		return 0, err
	}
	session.Commit()
	return email.Id, err
}
