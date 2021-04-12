package treasure_box

import (
	"game_server/core/utils"
)

type Base struct {
}

// generateResponseMessage 生成返回返回给客户端格式消息.
func (b *Base) ResponseMessage(opCode uint16, message interface{}) *utils.Packet {
	rsp := &utils.Packet{}
	rsp.Initialize(opCode)
	rsp.WriteData(message)
	return rsp
}

// NewBase 实例化Base结构体.
func NewBase() *Base {
	return &Base{}
}
