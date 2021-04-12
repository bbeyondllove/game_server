package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

type Search struct {
	b *Base
}

// ToSearchWord 获取热门和历史搜索词.
func (s *Search) TopSearchAndRecordWord(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_TOP_SEARCH_RESP, []string{})
	if packet != nil {
		return packet
	}

	// 请求本地生活服务接口.
	data, err := s.b.Request(LocalLifeApiSearchHistory, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_TOP_SEARCH_RESP, s.b.resMsg)
	}

	// todo 这里可以作其它业务处理.

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_TOP_SEARCH_RESP, s.b.resMsg)
}

// DeleteSearchRecord 删除历史搜索记录.
func (s *Search) DeleteSearchRecord(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD_RESP, []string{"name"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiDelSearchHistory, "POST", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD_RESP, s.b.resMsg)
}

// CityList 获取城市列表.
func (s *Search) CityList(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CITY_LIST_RESP, []string{})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiCityList, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_LIST_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_LIST_RESP, s.b.resMsg)
}

// StoreClassify 店铺分类信息.
func (s *Search) StoreClassify(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_STORE_ClASSIFY_RESP, []string{"category"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiStoreType, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_ClASSIFY_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_STORE_ClASSIFY_RESP, s.b.resMsg)
}

// categorySearch 分类搜索-店铺.
func (s *Search) CategoryStoreSearch(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE_RESP, []string{"searchType", "kind"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiCategorySearch, "POST", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE_RESP, s.b.resMsg)
}

// CategoryGoodsSearch 分类搜索-商品.
func (s *Search) CategoryGoodsSearch(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS_RESP, []string{"searchType", "kind"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiCategorySearch, "POST", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS_RESP, s.b.resMsg)
}

// SearchSuggest 搜索联想.
func (s *Search) SearchSuggest(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_SEARCH_SUGGEST_RESP, []string{"kind", "keyword"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiSearchSuggest, "POST", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_SEARCH_SUGGEST_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_SEARCH_SUGGEST_RESP, s.b.resMsg)
}

// TopSearchV2 历史和热门搜索V2版本.
func (s *Search) TopSearchV2(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_TOP_SEARCH_V2_RESP, []string{"kind"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiTopSearchV2, "GET", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_TOP_SEARCH_V2_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_TOP_SEARCH_V2_RESP, s.b.resMsg)
}

// DeleteSearchRecordV2 v2版本删除搜索历史记录.
func (s *Search) DeleteSearchRecordV2(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := s.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE_RESP, []string{"kind", "name"})
	if packet != nil {
		return packet
	}

	data, err := s.b.Request(LocalLifeApiSearchDelete, "DELETE", requestArgs)
	if err != nil {
		s.b.resMsg.Code = LocalLifeServerBusy
		s.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE_RESP, s.b.resMsg)
	}

	s.b.resMsg.Code = LocalLifeSuccess
	s.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	s.b.resMsg.Data = data

	return s.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE_RESP, s.b.resMsg)
}

// NewSearch 实例化Search结体体.
func NewSearch() *Search {
	return &Search{b: NewBase()}
}
