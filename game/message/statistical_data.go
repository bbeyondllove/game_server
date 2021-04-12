package message

import (
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/message/statistical_data"
)

// 不需要返回
func (a *agent) IconStatistics(requestMsg *utils.Packet) {
	logger.Infof("IconStatistics in request:", requestMsg.GetBuffer())
	data := statistical_data.NewDate()
	// 统计数据
	data.CityIcon(requestMsg)
	return
}

// 宝箱有效广告统计
func (a *agent) TreasureBoxStatistics(requestMsg *utils.Packet) {
	logger.Infof("TreasureBoxStatistics in request:", requestMsg.GetBuffer())
	data := statistical_data.NewDate()
	data.EffectiveAdvertis(requestMsg)
	logger.Debugf("TreasureBoxStatistics end")
	return
}
