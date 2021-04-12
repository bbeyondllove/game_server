package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"game_server/core/logger"
	"game_server/core/utils"
	"net"
	"runtime/debug"
	"sync/atomic"
	"time"

	"game_server/core/auth"
	"game_server/core/common"

	"github.com/gorilla/websocket"
)

type WSConn struct {
	conn       *websocket.Conn
	writeChan  *utils.SyncQueue
	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func NewWSConn(conn *websocket.Conn, maxMsgLen uint32) *WSConn {
	w := new(WSConn)
	w.conn = conn
	w.writeChan = utils.NewSyncQueue()
	w.maxMsgLen = maxMsgLen

	go func() {
		for {
			bs := w.writeChan.PopAll()
			if bs == nil {
				break
			}
			for _, b := range bs {
				buf, ok := b.([]byte)
				if ok {
					err := conn.WriteMessage(websocket.BinaryMessage, buf)
					if err != nil {
						logger.Info(err)
						goto closeSocket
					}
				}
			}

		}
	closeSocket:
		atomic.StoreInt32(&w.closeFlag, 1)
		conn.Close()
	}()

	return w
}

func (w *WSConn) Close() {
	logger.Errorf("server is closing conn, server:%s\n", string(debug.Stack()))
	w.writeChan.Close()
}

func (w *WSConn) IsClosed() bool {
	return atomic.LoadInt32(&w.closeFlag) == 1
}

func (w *WSConn) Write(b []byte) {
	w.writeChan.Push(b)
}

//不进队列，直接发送
func (w *WSConn) DirectWrite(b []byte) {
	w.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (w *WSConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WSConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WSConn) Read(b []byte) (n int, err error) {
	var l int
	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		return l, err
	}
	copy(b, msg)
	l = len(msg)
	if l < 0 {
		logger.Infof("ReadMessage error n=", n)
		return 0, errors.New("read empty")
	}
	w.activeTime = time.Now().Unix()
	// base.Log.Debug("w.activeTime=", w.activeTime)
	return l, nil
}

func (w *WSConn) ReadMsg() (*common.WorldPacket, error) {
	_, b, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	w.Crypt.DecryptRecv(b[:4])
	msgLen := int(binary.LittleEndian.Uint16(b[:2]))
	opCode := binary.LittleEndian.Uint16(b[2:4])
	if msgLen != len(b) {
		return nil, errors.New("收到ws数据长度错误")
	}
	packet := &common.WorldPacket{}
	packet.Initialize(opCode)
	packet.WriteBytes(b[4:])
	return packet, err
}

// args must not be modified by the others goroutines
func (w *WSConn) WriteMsg(packet *common.WorldPacket) error {
	if w.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 4)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	w.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	w.Write(header.Bytes())
	return nil
}

func (w *WSConn) InitCrypt(k []byte) {
	w.Crypt.Init(k)
}

func (w *WSConn) IsTimeout(maxTime uint32) bool {
	var ret bool
	now := time.Now().Unix()
	cur := uint32(now - w.activeTime)
	if cur > maxTime {
		logger.Debugf("IsTimeout() now=，,cur=,w.activeTime=", now, cur, w.activeTime)
		ret = true
	}
	return ret
}
