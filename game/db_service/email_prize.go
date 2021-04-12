package db_service

import (
	"game_server/db"
	"game_server/game/model"
)

type EmailPrizeInfo struct{}

const (
	EMAIL_PRIZE_TABLE = "t_email_prize"
)

/*
查询邮件奖品
*/
func (this *EmailPrizeInfo) GetEmailPrize(emailId int64) ([]model.EmailPrize, error) {
	emailPrizeList := make([]model.EmailPrize, 0)
	session := db.Mysql.Table(EMAIL_PRIZE_TABLE)
	session.Where("email_id = ?", emailId)
	session.Where("display = ?", 1)
	session.OrderBy("id desc")
	err := session.Find(&emailPrizeList)
	return emailPrizeList, err
}
