package message

import (
	"game_server/core/logger"
	"sync"

	"golang.org/x/net/websocket"
)

var connMgr = newConnMgr()

type ConnMgr struct {
	conns map[string]*websocket.Conn
	m     sync.Mutex
}

func newConnMgr() *ConnMgr {
	c := &ConnMgr{
		conns: make(map[string]*websocket.Conn),
	}
	return c
}

//移除连接
func (mgr *ConnMgr) CloseConn(userId string) {
	mgr.m.Lock()
	defer mgr.m.Unlock()

	for k, v := range mgr.conns {
		if k == userId {
			(*v).Close()
			logger.Infof("user exit:", k)
			delete(mgr.conns, k)
			break
		}
	}

}
