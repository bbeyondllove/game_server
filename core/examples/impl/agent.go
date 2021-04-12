package impl

import (
	"bytes"
	"encoding/binary"
	"kk_server/base"
	"kk_server/protocol"
	"net"

	"game_server/core/common"
	"game_server/core/network"
	"game_server/corelogger"

	"github.com/golang/protobuf/proto"
)

type agent struct {
	conn    network.Conn
	session *CSession
	auth    bool
}

//SendPacket send msg
func SendPacket(conn network.Conn, msg common.IPacket) error {

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, msg.Len())
	binary.Write(data, binary.LittleEndian, msg.GetCmd())
	binary.Write(data, binary.LittleEndian, msg.Bytes())
	conn.Write(data.Bytes())

	conn.Write(data.Bytes())
	return nil
}

//Run 在socket的recv gorouting执行
func (a *agent) Run() {
	logger.Debug("agent.Run()cstart")
	defer utils.SysRecoverWrap(Run)
	
	for {
		packet, err := a.conn.ReadMsg()
		if err != nil {
			logger.Error("err:", err.Error())
			a.Close()
			break
		}

		if packet.GetCmd() != uint16(protocol.Cmd_CBeat) {
			logger.Debug("packet.GetCmd()=", packet.GetCmd())
		}

		agentHandle := AgentCodeTable[packet.GetCmd()] //未登录业务处理
		if agentHandle.Handler != nil {
			agentHandle.Handler(a, packet)
		} else {

			if a.session != nil && a.auth {
				a.session.QueuePacket(packet)
			} else {
				logger.Error("session is nil")
				a.conn.Close()
				return
			}
		}

	}
}

func (a *agent) OnClose() {
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

//用户心跳
func (a *agent) HandleHEARTBEAT(requestMsg common.IPacket) {
	logger.Debug("Handle_HEARTBEAT in")

	rsp := &common.Packet{}
	rsp.Initialize(uint16(protocol.Cmd_SBeat))
	SendPacket(a.conn, rsp)

}

func (a *agent) HandleLogin(reqMsg common.IPacket) {
	logger.Debug("HandleLogin in")

	pbData := &protocol.ClientLogin{}
	err := proto.Unmarshal(reqMsg.Bytes(), pbData)
	if err != nil {
		logger.Error("proto.Unmarshal failed cmd=", reqMsg.GetCmd())
		return // 跳出循环，进行下一次消息读取
	}

}
