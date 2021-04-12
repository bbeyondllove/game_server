package message

import (
	"game_server/core/logger"
	"net/http"
	"sync"

	"game_server/core/network"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	maxConnNum int
	maxMsgLen  uint32
	newAgent   func(network.Conn) network.Agent
	upgrader   websocket.Upgrader
	conns      network.WsConnSet
	mt         sync.Mutex
	wg         sync.WaitGroup
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	protocol := r.Header.Get("Sec-WebSocket-Protocol")
	var header http.Header
	if protocol != "" {
		header = make(http.Header)
		header["Sec-WebSocket-Protocol"] = []string{protocol}
	}
	conn, err := h.upgrader.Upgrade(w, r, header)
	if err != nil {
		logger.Warnf("upgrade error: ", err)
		return
	}
	conn.SetReadLimit(int64(h.maxMsgLen))

	h.wg.Add(1)
	defer h.wg.Done()

	h.mt.Lock()
	if h.conns == nil {
		h.mt.Unlock()
		conn.Close()
		return
	}
	if len(h.conns) >= h.maxConnNum {
		h.mt.Unlock()
		conn.Close()
		logger.Warnf("too many connections")
		return
	}
	h.conns[conn] = struct{}{}
	h.mt.Unlock()

	wsConn := network.NewWSConn(conn, h.maxMsgLen)
	agent := h.newAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	h.mt.Lock()
	delete(h.conns, conn)
	h.mt.Unlock()
	agent.OnClose()
}

func (h *WSHandler) Close() {
	h.mt.Lock()
	for conn := range h.conns {
		conn.Close()
	}
	h.conns = nil
	h.mt.Unlock()

	h.wg.Done()
}
