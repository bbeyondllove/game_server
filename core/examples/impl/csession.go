package impl

import (
	"kk_server/base"
	"sync"
	"time"

	"game_server/core/common"
	"game_server/core/logger"
	"game_server/core/network"
	"game_server/core/utils"
)

const MaxTime = 2 * 60 //保活最大间隔2分钟

type CSession struct {
	UserId     int64
	recvQueue  *utils.SyncQueue
	conn       network.Conn
	logoutTime int64
	once       sync.Once
}

func NewCSession(userId int64, conn network.Conn) *CSession {
	cs := &CSession{
		UserId:    userId,
		recvQueue: utils.NewSyncQueue(),
		conn:      conn,
	}
	return cs
}

func (s *CSession) sendPacket(msg common.IPacket) {
	if SendPacket(s.conn, msg) != nil {
		s.conn.Close()
	}
}

func (s *CSession) Update() bool {
	pcks, ok := s.recvQueue.TryPopAll()
	if !ok || pcks == nil {
		goto check
	}
	for _, pck := range pcks {
		msg, ok := pck.(*common.Packet)
		if !ok {
			logger.Info("pck.(*proto.C2SMessage) not ok")
			break
		}
		s.handler(msg)
	}
check:
	if s.conn.IsClosed() && s.logoutTime == 0 { //socket is closed
		logger.Debug("should be closed()111")
		s.logoutTime = time.Now().Unix()
	}

	///- If necessary, logout
	currTime := time.Now().Unix()
	if s.ShouldLogOut(currTime) {
		logger.Debug("should be closed()222")
		return false // Will remove this session from session map
	}

	return true
}

func (s *CSession) ShouldLogOut(curTime int64) bool {
	var ret bool
	if s.logoutTime > 0 && curTime >= s.logoutTime {
		logger.Debugf("ShouldLogOut() curTime=,s.logoutTime=", curTime, s.logoutTime)
		ret = true
	}
	if s.conn.IsTimeout(MaxTime) {
		s.Close()
		logger.Debug("ShouldLogOut() s.conn.IsTimeout")
		ret = true
	}
	return ret
}

//延迟关闭session
func (s *CSession) CloseSession() {
	//防止出现前面的消息还没回复就断链了.
	timer := time.NewTimer(time.Second * 2)
	for {
		select {
		case <-timer.C:
			s.Close()
			return
		}
	}
}

//连接关闭
func (s *CSession) Close() {
	s.once.Do(func() {
		s.conn.Close()
	})
}

func (s *CSession) handler(msg common.IPacket) {
	defer utils.SysRecoverWrap(handler)

	logger.Debugf("CSession handler() msg.GetCmd()=,userid=", msg.GetCmd(), s.UserId)
	opHandle := OpCodeTable[msg.GetCmd()]
	if opHandle.Handler != nil {
		opHandle.Handler(s, msg)
	} else {
		logger.Errorf("unknown opcode:", msg.GetCmd())
	}
}

//QueuePacket 消息入队
func (s *CSession) QueuePacket(msg common.IPacket) {
	s.recvQueue.Push(msg)
}

func (s *CSession) HandleGetAmount(packet common.IPacket) {

}
