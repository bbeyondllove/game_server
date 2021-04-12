package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

// Hotel 酒店结构体.
type Hotel struct {
	b *Base
}

// Search 搜索酒店.
func (h *Hotel) Search(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := h.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_HOTEL_SEARCH_RESP, []string{})
	if packet != nil {
		return packet
	}

	// 请求本地生活服务接口.
	data, err := h.b.Request(LocalLifeApiHotelSearch, "POST", requestArgs)
	if err != nil {
		h.b.resMsg.Code = LocalLifeServerBusy
		h.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_HOTEL_SEARCH_RESP, h.b.resMsg)
	}

	// todo 这里可以作其它业务处理.

	h.b.resMsg.Code = LocalLifeSuccess
	h.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	h.b.resMsg.Data = data

	return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_HOTEL_SEARCH_RESP, h.b.resMsg)
}

// GetDetails　获取酒店详情.
func (h *Hotel) GetDetails(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := h.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_HOTEL_DETAIL_RESP, []string{"storeId"})
	if packet != nil {
		return packet
	}

	// 请求本地生活服务接口.
	data, err := h.b.Request(LocalLifeApiHotelDetail, "GET", requestArgs)
	if err != nil {
		h.b.resMsg.Code = LocalLifeServerBusy
		h.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_HOTEL_DETAIL_RESP, h.b.resMsg)
	}

	h.b.resMsg.Code = LocalLifeSuccess
	h.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	h.b.resMsg.Data = data

	return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_HOTEL_DETAIL_RESP, h.b.resMsg)
}

// GetRoomDetails 获取房间详情.
func (h *Hotel) GetRoomDetails(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := h.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_ROOM_DETAIL_RESP, []string{"goodsId"})
	if packet != nil {
		return packet
	}

	// 请求本地生活服务接口.
	data, err := h.b.Request(LocalLifeApiHotelRoomDetail, "GET", requestArgs)
	if err != nil {
		h.b.resMsg.Code = LocalLifeServerBusy
		h.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_ROOM_DETAIL_RESP, h.b.resMsg)
	}

	h.b.resMsg.Code = LocalLifeSuccess
	h.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	h.b.resMsg.Data = data

	return h.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_ROOM_DETAIL_RESP, h.b.resMsg)
}

// NewHotel 实例化Hotel结体体.
func NewHotel() *Hotel {
	return &Hotel{b: NewBase()}
}
