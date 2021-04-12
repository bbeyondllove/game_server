package local_life

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

type City struct {
	b *Base
}

// CityInfos 查询城市列表.
func (c *City) CityInfos(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := c.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CITY_RESP, []string{})
	if packet != nil {
		return packet
	}

	data, err := c.b.Request(LocalLifeApiCity, "GET", requestArgs)
	if err != nil {
		c.b.resMsg.Code = LocalLifeServerBusy
		c.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_RESP, c.b.resMsg)
	}

	c.b.resMsg.Code = LocalLifeSuccess
	c.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	c.b.resMsg.Data = data

	return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_RESP, c.b.resMsg)
}

// CitySuggest 城市联想.
func (c *City) CitySuggest(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := c.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CITY_SUGGEST_RESP, []string{"cityName"})
	if packet != nil {
		return packet
	}

	data, err := c.b.Request(LocalLifeApiCitySuggest, "POST", requestArgs)
	if err != nil {
		c.b.resMsg.Code = LocalLifeServerBusy
		c.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_SUGGEST_RESP, c.b.resMsg)
	}

	c.b.resMsg.Code = LocalLifeSuccess
	c.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	c.b.resMsg.Data = data

	return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_SUGGEST_RESP, c.b.resMsg)
}

// CityPick 切换城市.
func (c *City) CityPick(requestMsg *utils.Packet) *utils.Packet {
	requestArgs, packet := c.b.beforeRequest(requestMsg, proto.MSG_LOCAL_LIFE_CITY_PICK_RESP, []string{"cityName"})
	if packet != nil {
		return packet
	}

	data, err := c.b.Request(LocalLifeApiCityPick, "POST", requestArgs)
	if err != nil {
		c.b.resMsg.Code = LocalLifeServerBusy
		c.b.resMsg.Msg = statusCodeMessage[LocalLifeServerBusy]
		return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_PICK_RESP, c.b.resMsg)
	}

	c.b.resMsg.Code = LocalLifeSuccess
	c.b.resMsg.Msg = statusCodeMessage[LocalLifeSuccess]
	c.b.resMsg.Data = data

	return c.b.generateResponseMessage(proto.MSG_LOCAL_LIFE_CITY_PICK_RESP, c.b.resMsg)
}

// NewCity 实例化City.
func NewCity() *City {
	return &City{b: NewBase()}
}
