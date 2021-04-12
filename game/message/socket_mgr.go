package message

import (
	"net/http"
	"time"

	"game_server/core/base"

	"game_server/core/network"

	"github.com/gorilla/websocket"
)

type SocketMgr struct {
	wsServer *network.WSServer
}

func (mgr *SocketMgr) Init() {

	if base.Setting.Server.LocalHost != "" {
		mgr.wsServer = new(network.WSServer)
		mgr.wsServer.Addr = base.Setting.Server.LocalHost + ":" + base.Setting.Server.LocalPort
		mgr.wsServer.MaxConnNum = base.Setting.Server.MaxConnNum
		mgr.wsServer.MaxMsgLen = base.Setting.Server.MaxMsgLen
		mgr.wsServer.HTTPTimeout = base.Setting.Server.HttpTimeout * time.Second
		mgr.wsServer.CertFile = base.Setting.Server.CertFile
		mgr.wsServer.KeyFile = base.Setting.Server.KeyFile
		mgr.wsServer.NewAgent = func(conn network.Conn) network.Agent {
			a := &agent{conn: conn}
			return a
		}
		mgr.wsServer.Handler = &WSHandler{
			maxConnNum: base.Setting.Server.MaxConnNum,
			maxMsgLen:  base.Setting.Server.MaxMsgLen,
			newAgent:   mgr.wsServer.NewAgent,
			conns:      make(network.WsConnSet),
			upgrader: websocket.Upgrader{
				HandshakeTimeout: base.Setting.Server.HttpTimeout * time.Second,
				CheckOrigin:      func(_ *http.Request) bool { return true },
				//Subprotocols:     []string{"binary"},
			},
		}
	}

	if mgr.wsServer != nil {
		mgr.wsServer.Start()
	}
}

func (mgr *SocketMgr) Destroy() {
	if mgr.wsServer != nil {
		mgr.wsServer.Close()
	}
}
