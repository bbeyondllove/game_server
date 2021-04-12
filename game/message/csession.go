package message

import (
	"game_server/core/utils"
	"runtime/debug"
	"sync"
	"time"

	"game_server/core/logger"
	"game_server/core/network"
)

const MaxTime = 1 * 60 //保活最大间隔1分钟

type CSession struct {
	UserId     string
	LocateId   int32
	X          int    //坐标x
	Y          int    //坐标
	Platform   int    //平台
	Version    string //平台版本
	recvQueue  *utils.SyncQueue
	conn       network.Conn
	logoutTime int64
	once       sync.Once
	BReconnect bool
}

func NewCSession(locationId int32, userId string, platform int, version string, conn network.Conn) *CSession {
	cs := &CSession{
		UserId:     userId,
		LocateId:   locationId,
		Platform:   platform,
		Version:    version,
		recvQueue:  utils.NewSyncQueue(),
		conn:       conn,
		BReconnect: false,
	}
	return cs
}

func (s *CSession) sendPacket(msg *utils.Packet) {

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
		msg, ok := pck.(*utils.Packet)
		if !ok {
			logger.Infof("pck.(*proto.C2SMessage) not ok")
			break
		}
		s.handler(msg)
	}
check:
	if s.conn.IsClosed() && s.logoutTime == 0 { //socket is closed
		logger.Debugf("should be closed()111")
		s.logoutTime = time.Now().Unix()
	}

	///- If necessary, logout
	currTime := time.Now().Unix()
	if s.ShouldLogOut(currTime) {
		logger.Debugf("should be closed()222")
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
		logger.Debugf("ShouldLogOut() s.conn.IsTimeout")
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

func (s *CSession) handler(msg *utils.Packet) {
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("CSession panic err:[%V] ", string(debug.Stack()))
		}
	}()
	opHandle := OpCodeTable[msg.OpCode]
	if opHandle.Handler != nil {
		logger.Debugf("CSession handler() msg.MsgType=,userid=", opHandle.Name, s.UserId)
		opHandle.Handler(s, msg)
	} else {
		logger.Errorf("unknown opcode:", msg.OpCode)
	}
}

//QueuePacket 消息入队
func (s *CSession) QueuePacket(msg *utils.Packet) {
	s.recvQueue.Push(msg)
}
