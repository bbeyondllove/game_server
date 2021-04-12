package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"strconv"
	"time"

	"game_server/core/logger"

	"github.com/go-xorm/xorm"
)

type GameWallet struct{}

const (
	wallet_nid  = 1
	WalletTable = "t_game_wallet"
)

//获取指定金额
func (gw *GameWallet) GetAmount(user_id string, token_code string) (*model.GameWallet, error) {
	data := &model.GameWallet{}
	_, err := db.Mysql.Table(WalletTable).Where("user_id= ? and token_code=?", user_id, token_code).Get(data)
	return data, err
}

//获取所有金额
func (gw *GameWallet) GetAllAmount(user_id string) ([]model.GameWallet, error) {
	data := []model.GameWallet{}
	err := db.Mysql.Table(WalletTable).Where("user_id= ?", user_id).Find(&data)
	return data, err
}

//获取自动添加序列号
func (gw *GameWallet) NextSeqId(session *xorm.Session) (string, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.GameWallet{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return "", err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return strconv.Itoa(maxId + 1), nil
	} else {
		count, err := db.Mysql.Count(&model.GameWallet{})
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(count + wallet_nid)), nil
	}
}

//添加记录
func (gw *GameWallet) Add(data_map *model.GameWallet) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	seqId, err := gw.NextSeqId(session)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}

	data_map.Id = seqId
	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}
	return true, nil
}

//更新金额
func (u *GameWallet) UpdateMoney(userId, tokenCode string, money float64) (bool, *xorm.Session) {
	var sql string
	session := db.Mysql.NewSession()
	session.Begin()

	moneyStr := strconv.FormatFloat(money, 'E', -1, 64)
	sql = "update " + WalletTable + " SET	amount = amount+ " + moneyStr + ",amount_available = amount_available + " + moneyStr + " where user_id=" + userId + " and token_code='" + tokenCode + "'"
	r, err := db.Mysql.Exec(sql)
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: ", err)
		return false, session
	}

	return true, session
}
