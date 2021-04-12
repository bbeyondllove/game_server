package proto

import (
	"github.com/shopspring/decimal"
)

//购买商品
type C2SBuyItem struct {
	Token  string `json:"token"`
	ItemId int    `json:"itemId"` //商品ID
}

type S2CBuyItem struct {
	S2CCommon
	ItemId ProductItem     `json:"itemId"` //商品ID
	Money  decimal.Decimal `json:"money"`
}

//获取商品列表
type C2SGetItemList struct {
	Token string `json:"token"`
}
