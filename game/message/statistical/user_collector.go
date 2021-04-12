package statistical

import (
	"game_server/db"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
)

type UserDataCollector struct {
}

// 获取用户在线时长
func (this *UserDataCollector) GetOnlineTime(userId string) int {
	key := proto.ONLINE_KEY + userId
	valueRet, err := db.RedisGame.HGet(key, "static").Result()
	if err != nil {
		return 0
	}
	online_time, err := strconv.Atoi(valueRet)
	if err != nil {
		return 0
	}
	return online_time
}

// 用户当天产出CDT总数
func (this *UserDataCollector) GetCDTCount(userId, day string) float32 {
	num, err := model.NewCdtRecord().GetUserOneDayCdt(userId, day, 0)
	if err != nil {
		return 0
	}
	return num
}
