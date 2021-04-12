package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

// Store 店铺struct.
type Store struct {
	b *Base
}

// StoreType 查询店铺分类信息.
func (s *Store) StoreType(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_STORE_TYPE_RESP, []string{"kind"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiStoreTypeV2, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_TYPE_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_TYPE_RESP, s.b.resMsg)
}

// StoreHotelDetail 住宿详情.
func (s *Store) StoreHotelDetail(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL_RESP, []string{"kind", "storeId"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiStoreHotelDetail, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL_RESP, s.b.resMsg)
}

// StoreRestaurantDetail 美食详情.
func (s *Store) StoreRestaurantDetail(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL_RESP, []string{"kind", "storeId"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiStoreRestaurantDetail, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL_RESP, s.b.resMsg)
}

// NewStore 实例化Store.
func NewStore() *Store {
	return &Store{b: NewBase()}
}
