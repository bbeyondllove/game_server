package statistical_data

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/proto"
	"time"
)

type Data struct {
}

func NewDate() *Data {
	return &Data{}
}

//统计城市中心ICON入口统计  pv, uv
func (d *Data) CityIcon(requestMsg *utils.Packet) {
	msg := &proto.C2SCityIcon{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("CityIcon json.Unmarshal error, err=", err.Error())
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("CityIcon token error", msg.Token)
		return
	}
	//获取当前时间 2006-01-02 15:04:05
	date := time.Now().Format("20060102")
	pvKey := CITY_ICON_PV + date
	uvKey := CITY_ICON_UV + date
	// pv
	isExists := db.RedisMgr.KeyExist(pvKey)
	if !isExists {
		db.RedisMgr.Incr(pvKey)
		db.RedisMgr.Expire(pvKey, 3600*24*3) // 过期时间设置3两天
	} else {
		db.RedisMgr.Incr(pvKey)
	}

	//uv
	isExists = db.RedisMgr.KeyExist(uvKey)
	if !isExists {
		db.RedisMgr.HIncrBy(uvKey, payLoad.UserId, 1)
		db.RedisMgr.Expire(uvKey, 3600*24*3) // 过期时间设置3两天
	} else {
		db.RedisMgr.HIncrBy(uvKey, payLoad.UserId, 1)
	}

}

// 统计有效广告
func (d *Data) EffectiveAdvertis(requestMsg *utils.Packet) {
	msg := &proto.C2SCEffectiveAdvertis{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("EffectiveAdvertis json.Unmarshal error, err=", err.Error())
		return
	}
	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("EffectiveAdvertis token error", msg.Token)
		return
	}
	requestId := msg.RequestID
	if requestId != "" {
		//统计总的
		date := time.Now().Format("20060102")
		key := BAOX_ADVERTIS + date
		isExists := db.RedisMgr.KeyExist(key)
		if !isExists {
			db.RedisMgr.HIncrBy(key, requestId, 1)
			db.RedisMgr.Expire(key, 3600*24*3) // 过期时间设置3两天
		} else {
			db.RedisMgr.HIncrBy(key, requestId, 1)
		}
		//统计每个用户的
		userKey := BAOX_ADVERTIS_USER + date + ":" + payLoad.UserId
		isExists = db.RedisMgr.KeyExist(userKey)
		if !isExists {
			db.RedisMgr.HIncrBy(userKey, requestId, 1)
			db.RedisMgr.Expire(userKey, 3600*24*3) // 过期时间设置3两天
		} else {
			db.RedisMgr.HIncrBy(userKey, requestId, 1)
		}
	}
}
