package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

// Index struct.
type Index struct {
	b *Base
}

// Recommend 首页推荐.
func (i *Index) Recommend(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := i.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_INDEX_RECOMMEND_RESP, []string{})
	if packet != nil {
		return packet
	}

	// 请求本地生活服务接口.
	data, err := i.b.Request(LocalLifeIndexRecommend, "GET", requestArgs)
	if err != nil {
		i.b.resMsg.Code = LocalLifeServerBusy
		i.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return i.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_INDEX_RECOMMEND_RESP, i.b.resMsg)
	}

	// todo 这里可以作其它业务处理.

	i.b.resMsg.Code = LocalLifeSuccess
	i.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	i.b.resMsg.Data = data

	return i.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_INDEX_RECOMMEND_RESP, i.b.resMsg)
}

// SelectCity 首页选择城市，主要用于后端城市统计切换，调整运营策略.
func (i *Index) SelectCity(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := i.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_INDEX_SELECT_CITY_RESP, []string{"cityName"})
	if packet != nil {
		return packet
	}

	data, err := i.b.Request(LocalLifeApiIndexSelectCity, "POST", requestArgs)
	if err != nil {
		i.b.resMsg.Code = LocalLifeServerBusy
		i.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return i.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_INDEX_SELECT_CITY_RESP, i.b.resMsg)
	}

	i.b.resMsg.Code = LocalLifeSuccess
	i.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	i.b.resMsg.Data = data

	return i.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_INDEX_SELECT_CITY_RESP, i.b.resMsg)
}

// NewIndex 实例化Index结构体.
func NewIndex() *Index {
	return &Index{b: NewBase()}
}
