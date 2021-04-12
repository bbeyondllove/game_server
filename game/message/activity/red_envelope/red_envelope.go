package red_envelope

import (
	"game_server/core/logger"
	"game_server/game/db_service"
	"game_server/game/proto"
	"math/rand"
	"sync"
	"time"
)

var (
	G_StreasureBox_Red_Envelope_Table sync.Map
	RedEnvelopeRandTable              = "RedEnvelopeRandTable"
	Min                               float32 //最小值
	Max                               float32 //最大值
)

// 通过数据库初始化随机表
func InitRedEnvelopeRandTable(randTable string, randType int) {
	data, err := db_service.TreasureBoxCdtConfigIns.GetAllData(randType)
	if err != nil {
		logger.Errorf("InitRedEnvelopeRandTable error,err= %v", err)
		return
	}
	if len(data) == 0 {
		logger.Error("InitRedEnvelopeRandTable is empty")
		return
	}
	Min = data[0].RewardNumber
	Max = data[0].RewardNumber
	last_level := int64(0)
	var rand_table []proto.ReaEnvelRand = make([]proto.ReaEnvelRand, 0)
	for _, item := range data {
		// 最小值
		if Min > item.RewardNumber {
			Min = item.RewardNumber
		}
		// 最大值
		if Max < item.RewardNumber {
			Max = item.RewardNumber
		}
		level := item.Probability + last_level
		last_level = level
		rand_table = append(rand_table, proto.ReaEnvelRand{
			Probability:  level,
			RewardNumber: item.RewardNumber,
		})
	}
	G_StreasureBox_Red_Envelope_Table.Store(randTable, rand_table)
}

func RedEnvelopeConfig() {
	InitRedEnvelopeRandTable(RedEnvelopeRandTable, 2)
}

type RedEnvelope struct {
}

func NewRedEnvelope() *RedEnvelope {
	return &RedEnvelope{}
}

// 随机获取奖励
func (this *RedEnvelope) RandCdtAward(randTable string) float32 {
	table := GetRandTable(randTable)
	table_len := len(table)
	if table_len == 0 {
		return Min
	}
	total_len := table[table_len-1].Probability
	award := Min
	for {
		rand.Seed(time.Now().UTC().UnixNano())
		idx := rand.Int63n(total_len)
		logger.Debugf("RandCdtAward: %v", idx)
		for _, item := range table {
			if idx < item.Probability {
				award = item.RewardNumber
				break
			}
		}
		//不放出 一等奖
		if award < Max {
			break
		}
	}
	return award
}

func GetRandTable(randTable string) []proto.ReaEnvelRand {
	result, ok := G_StreasureBox_Red_Envelope_Table.Load(randTable)
	if ok {
		return result.([]proto.ReaEnvelRand)
	}
	return nil
}
