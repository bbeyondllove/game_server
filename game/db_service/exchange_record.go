package db_service

import (
	"game_server/core/base"
	"game_server/db"
	"game_server/game/model"

	"game_server/core/logger"
)

//交易记录
type ExchangeRecord struct{}

const exchange_table = "t_exchange_record"

//增加交易金额
func (gw *ExchangeRecord) Add(data_map *model.ExchangeRecord) (bool, error) {
	// data_map.CreateTime = time.Now()
	// data_map.UpdateTime = time.Now()

	_, err := db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: %v", err)
		return false, err
	}

	if err != nil {
		logger.Errorf("[Error]: %v", err)
		return false, err
	}
	return true, nil

}

func (gw *ExchangeRecord) GetExchangeList(uid string, lastRecordId, deadLine int64) ([]model.ExchangeRecord, error) {
	data := []model.ExchangeRecord{}
	var err error
	if deadLine > 0 {
		err = db.Mysql.Table(exchange_table).Where("id > ? and user_id= ?  and create_time <? limit ?", lastRecordId, uid, deadLine, base.Setting.Base.GetExchangeListMax).Find(data)
	} else {
		err = db.Mysql.Table(exchange_table).Where("id > ? and user_id= ? limit ?", lastRecordId, uid, base.Setting.Base.GetExchangeListMax).Find(data)
	}

	return data, err
}

func (gw *ExchangeRecord) GetExchangeByTxid(txid string) (*model.ExchangeRecord, error) {
	data := &model.ExchangeRecord{}
	_, err := db.Mysql.Table(exchange_table).Where("sys_order_sn= ?", txid).Get(data)
	return data, err

}
