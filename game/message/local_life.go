package message

import (
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/message/local_life"
)

// LocalLife 本地生活struct.
type LocalLife struct {
}

// Recommend 首页推荐.
func (a *agent) Recommend(requestMsg *utils.Packet) {
	logger.Infof("Recommend in request:", requestMsg.GetBuffer())
	index := local_life.NewIndex()
	rsp := index.Recommend(requestMsg)
	SendPacket(a.conn, rsp)
}

// TopSearchAndRecordWord 获取热门和历史搜索词.
func (a *agent) TopSearchAndRecordWord(requestMsg *utils.Packet) {
	logger.Infof("TopSearchAndRecordWord in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.TopSearchAndRecordWord(requestMsg)
	SendPacket(a.conn, rsp)
}

// DeleteSearchRecord 删除历史搜索记录.
func (a *agent) DeleteSearchRecord(requestMsg *utils.Packet) {
	logger.Infof("DeleteSearchRecord in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.DeleteSearchRecord(requestMsg)
	SendPacket(a.conn, rsp)
}

// CityList 获取城市列表.
func (a *agent) CityList(requestMsg *utils.Packet) {
	logger.Infof("CityList in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.CityList(requestMsg)
	SendPacket(a.conn, rsp)
}

// StoreClassify 店铺分类信息.
func (a *agent) StoreClassify(requestMsg *utils.Packet) {
	logger.Infof("StoreClassify in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.StoreClassify(requestMsg)
	SendPacket(a.conn, rsp)
}

// SearchHotel 搜索酒店.
func (a *agent) SearchHotel(requestMsg *utils.Packet) {
	logger.Infof("SearchHotel in request:", requestMsg.GetBuffer())
	hotel := local_life.NewHotel()
	rsp := hotel.Search(requestMsg)
	SendPacket(a.conn, rsp)
}

// GetHotelDetails 获取酒店详情.
func (a *agent) GetHotelDetails(requestMsg *utils.Packet) {
	logger.Infof("GetHotelDetails in request:", requestMsg.GetBuffer())
	hotel := local_life.NewHotel()
	rsp := hotel.GetDetails(requestMsg)
	SendPacket(a.conn, rsp)
}

// GetHotelRoomDetails 获取酒店房间详情.
func (a *agent) GetHotelRoomDetails(requestMsg *utils.Packet) {
	logger.Infof("GetHotelRoomDetails in request:", requestMsg.GetBuffer())
	hotel := local_life.NewHotel()
	rsp := hotel.GetRoomDetails(requestMsg)
	SendPacket(a.conn, rsp)
}

// IndexSelectCity 首页选择城市，主要用于后端城市统计切换，调整运营策略.
func (a *agent) IndexSelectCity(requestMsg *utils.Packet) {
	logger.Infof("IndexSelectCity in request:", requestMsg.GetBuffer())
	index := local_life.NewIndex()
	rsp := index.SelectCity(requestMsg)
	SendPacket(a.conn, rsp)
}

// CategoryStoreSearch 分类搜索-店铺.
func (a *agent) CategoryStoreSearch(requestMsg *utils.Packet) {
	logger.Infof("CategoryStoreSearch in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.CategoryStoreSearch(requestMsg)
	SendPacket(a.conn, rsp)
}

// CategoryGoodsSearch 分类搜索-商品.
func (a *agent) CategoryGoodsSearch(requestMsg *utils.Packet) {
	logger.Infof("CategoryGoodsSearch in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.CategoryGoodsSearch(requestMsg)
	SendPacket(a.conn, rsp)
}

// SearchSuggest 搜索联想.
func (a *agent) SearchSuggest(requestMsg *utils.Packet) {
	logger.Infof("SearchSuggest in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.SearchSuggest(requestMsg)
	SendPacket(a.conn, rsp)
}

// TopSearchV2 搜索联想.
func (a *agent) TopSearchV2(requestMsg *utils.Packet) {
	logger.Infof("TopSearchV2 in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.TopSearchV2(requestMsg)
	SendPacket(a.conn, rsp)
}

// StoreType 查询店铺分类信息.
func (a *agent) StoreType(requestMsg *utils.Packet) {
	logger.Infof("StoreType in request:", requestMsg.GetBuffer())
	store := local_life.NewStore()
	rsp := store.StoreType(requestMsg)
	SendPacket(a.conn, rsp)
}

// StoreHotelDetail 住宿详情.
func (a *agent) StoreHotelDetail(requestMsg *utils.Packet) {
	logger.Infof("StoreHotelDetail in request:", requestMsg.GetBuffer())
	store := local_life.NewStore()
	rsp := store.StoreHotelDetail(requestMsg)
	SendPacket(a.conn, rsp)
}

// StoreRestaurantDetail 美食详情.
func (a *agent) StoreRestaurantDetail(requestMsg *utils.Packet) {
	logger.Infof("StoreRestaurantDetail in request:", requestMsg.GetBuffer())
	store := local_life.NewStore()
	rsp := store.StoreRestaurantDetail(requestMsg)
	SendPacket(a.conn, rsp)
}

// GoodsHotelDetail 商品住宿(房间详情).
func (a *agent) GoodsHotelDetail(requestMsg *utils.Packet) {
	logger.Infof("GoodsHotelDetail in request:", requestMsg.GetBuffer())
	goods := local_life.NewGoods()
	rsp := goods.GoodsHotelDetail(requestMsg)
	SendPacket(a.conn, rsp)
}

// DiscountDetail 美食优惠劵.
func (a *agent) DiscountDetail(requestMsg *utils.Packet) {
	logger.Infof("DiscountDetail in request:", requestMsg.GetBuffer())
	goods := local_life.NewGoods()
	rsp := goods.DiscountDetail(requestMsg)
	SendPacket(a.conn, rsp)
}

// RestaurantDetail 美食套餐.
func (a *agent) RestaurantDetail(requestMsg *utils.Packet) {
	logger.Infof("RestaurantDetail in request:", requestMsg.GetBuffer())
	goods := local_life.NewGoods()
	rsp := goods.RestaurantDetail(requestMsg)
	SendPacket(a.conn, rsp)
}

// CityInfos 查询城市列表.
func (a *agent) CityInfos(requestMsg *utils.Packet) {
	logger.Infof("CityInfos in request:", requestMsg.GetBuffer())
	city := local_life.NewCity()
	rsp := city.CityInfos(requestMsg)
	SendPacket(a.conn, rsp)
}

// CitySuggest 城市联想.
func (a *agent) CitySuggest(requestMsg *utils.Packet) {
	logger.Infof("CitySuggest in request:", requestMsg.GetBuffer())
	city := local_life.NewCity()
	rsp := city.CitySuggest(requestMsg)
	SendPacket(a.conn, rsp)
}

// CityPick 切换城市.
func (a *agent) CityPick(requestMsg *utils.Packet) {
	logger.Infof("CityPick in request:", requestMsg.GetBuffer())
	city := local_life.NewCity()
	rsp := city.CityPick(requestMsg)
	SendPacket(a.conn, rsp)
}

// SearchRecordDeleteV2 v2版本－删除搜索历史.
func (a *agent) SearchRecordDeleteV2(requestMsg *utils.Packet) {
	logger.Infof("SearchRecordDeleteV2 in request:", requestMsg.GetBuffer())
	search := local_life.NewSearch()
	rsp := search.DeleteSearchRecordV2(requestMsg)
	SendPacket(a.conn, rsp)
}
