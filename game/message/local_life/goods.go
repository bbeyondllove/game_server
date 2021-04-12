package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

// Goods 商品struct.
type Goods struct {
	b *Base
}

// GoodsHotel 商品－住宿(房间详情)/团购-套餐详情/代金劵－优惠详情.
func (g *Goods) GoodsHotelDetail(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := g.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL_RESP, []string{"goodsId", "kind"})
	if packet != nil {
		return packet
	}

	data, err := g.b.Request(LocalLifeApiGoodsHotelDetail, "GET", requestArgs)
	if err != nil {
		g.b.resMsg.Code = LocalLifeServerBusy
		g.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL_RESP, g.b.resMsg)
	}

	g.b.resMsg.Code = LocalLifeSuccess
	g.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	g.b.resMsg.Data = data

	return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL_RESP, g.b.resMsg)
}

// RestaurantDetail 美食套餐.
func (g *Goods) RestaurantDetail(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := g.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL_RESP, []string{"goodsId"})
	if packet != nil {
		return packet
	}

	data, err := g.b.Request(LocalLifeApiGoodsRestaurantDetail, "POST", requestArgs)
	if err != nil {
		g.b.resMsg.Code = LocalLifeServerBusy
		g.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL_RESP, g.b.resMsg)
	}

	g.b.resMsg.Code = LocalLifeSuccess
	g.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	g.b.resMsg.Data = data

	return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL_RESP, g.b.resMsg)
}

// DiscountDetail 美食优惠劵.
func (g *Goods) DiscountDetail(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := g.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_DISCOUNT_DETAIL_RESP, []string{"goodsId"})
	if packet != nil {
		return packet
	}

	data, err := g.b.Request(LocalLifeApiGoodsDiscountDetail, "POST", requestArgs)
	if err != nil {
		g.b.resMsg.Code = LocalLifeServerBusy
		g.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_DISCOUNT_DETAIL_RESP, g.b.resMsg)
	}

	g.b.resMsg.Code = LocalLifeSuccess
	g.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	g.b.resMsg.Data = data

	return g.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_DISCOUNT_DETAIL_RESP, g.b.resMsg)
}

// NewGoods 实例化Goods.
func NewGoods() *Goods {
	return &Goods{b: NewBase()}
}
