package impl

import (
	"kk_server/protocol"

	"game_server/core/common"
)

var OpCodeTable = make(map[uint16]OpCodeHandler, 0)
var AgentCodeTable = make(map[uint16]AgentHandler, 0)

const (
	STATUS_NOT_AUTHED = 0
	STATUS_AUTHED     = 1
	STATUS_UNHANDLED  = 2
)

//消息路由器
type OpCodeHandler struct {
	Name    string
	Status  uint8
	Handler func(*CSession, common.IPacket)
}

type AgentHandler struct {
	Name    string
	Status  uint8
	Handler func(*agent, common.IPacket)
}

func init() {
	// OpCodeTable[uint16(protocol.MSG_GET_AMOUNT)] = OpCodeHandler{"MSG_GET_AMOUNT", STATUS_AUTHED, (*CSession).HandleGetAmount}

	AgentCodeTable[uint16(protocol.Cmd_CBeat)] = AgentHandler{"Cmd_CBeat", STATUS_NOT_AUTHED, (*agent).HandleHEARTBEAT}
	AgentCodeTable[uint16(protocol.Cmd_CLogin)] = AgentHandler{"Cmd_CLogin", STATUS_NOT_AUTHED, (*agent).HandleLogin}

}
