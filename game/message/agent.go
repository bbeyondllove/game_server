package message

import (
	"bytes"
	"encoding/binary"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/game/proto"
	"io"
	"net"

	"game_server/core/logger"
	"game_server/core/network"
)

type agent struct {
	conn    network.Conn
	session *CSession
	auth    bool
	token   string
}

//SendPacket send msg
func SendPacket(conn network.Conn, msg *utils.Packet) error {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, msg.OpCode)
	binary.Write(data, binary.LittleEndian, msg.Bytes())
	conn.Write(data.Bytes())
	return nil
}

//Run 在socket的recv gorouting执行
func (a *agent) Run() {
	logger.Debug("agent start")

	defer utils.SysRecoverWrap(a.Run)

	var buff = make([]byte, base.Setting.Server.MaxMsgLen)

	for {
		n, err := a.conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				logger.Errorf("read EOF from connect")
				continue
			}
			logger.Errorf("read data from socket failed: %v, number byte:%v, token is: %v\n", err, n, a.token)
			//logger.Errorf("err=%+v,token=%+v", err.Error(), a.token)
			a.Close()
			break
		}

		cmd := binary.LittleEndian.Uint16(buff[:2])
		packet := &utils.Packet{}
		packet.Initialize(cmd)
		packet.WriteBytes(buff[2:n])
		//logger.Debug("msg.MsgType=", cmd)

		agentHandle := AgentCodeTable[cmd] //未登录业务处理
		if agentHandle.Handler != nil {
			if cmd != proto.MSG_HEARTBEAT {
				logger.Debugf("msg.MsgType=", agentHandle.Name)
			}
			agentHandle.Handler(a, packet)

		} else {
			if a.session != nil && a.auth {
				a.session.QueuePacket(packet)
			} else {
				logger.Errorf("session is nil")
				a.conn.Close()
				return
			}
		}

		copy(buff, buff[n:])
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
