package db_service

import "game_server/game/model"

type Action struct {
	Handler func(userId string, prize model.EmailPrize) (bool, error)
}

// 邮件领取奖励 （根据奖励类型来读取对应的方法）
var EmailPrizeAction = make(map[int]Action, 0)

func init() {
	EmailPrizeAction[1] = Action{EmailLogicIns.ReloadUserKnapsack} // 奖励发放到背包
	EmailPrizeAction[2] = Action{EmailLogicIns.ActivityReward}     // 双蛋活动 领取奖励
}
