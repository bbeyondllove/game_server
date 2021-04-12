package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"

	"log"

	"game_server/core/auth"
	"game_server/core/common"

	"game_server/core/utils"
)

//ConnSet server conn map
type ConnSet map[net.Conn]struct{}

const HeaderLen uint32 = 4

//TCPConn tcp连接类
type TCPConn struct {
	conn       net.Conn
	writeChan  *utils.SyncQueue
	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func newTCPConn(conn net.Conn, maxMsgLen uint32) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = utils.NewSyncQueue()
	tcpConn.maxMsgLen = maxMsgLen
	//tcp写goroutine
	go func() {
		var tmp bytes.Buffer
		for {
			time.Sleep(100 * time.Millisecond)
			bs := tcpConn.writeChan.PopAll()
			if bs == nil {
				break
			}
			for _, b := range bs {
				buf, ok := b.([]byte)
				if ok {
					tmp.Write(buf)
				}
			}
			if tmp.Len() != 0 {
				_, err := conn.Write(tmp.Bytes())
				if err != nil {
					log.Print(err.Error())
					goto closeSocket
				}
			}
			tmp.Reset()
		}
	closeSocket:
		atomic.StoreInt32(&tcpConn.closeFlag, 1)
		conn.Close()
	}()

	return tcpConn
}

func (tcpConn *TCPConn) IsClosed() bool {
	return atomic.LoadInt32(&tcpConn.closeFlag) == 1
}
func (tcpConn *TCPConn) Close() {
	tcpConn.writeChan.Close()
}

func (tcpConn *TCPConn) Write(b []byte) {
	tcpConn.writeChan.Push(b)
}

//不进队列，直接发送
func (tcpConn *TCPConn) DirectWrite(b []byte) {
	tcpConn.conn.Write(b)
}
func (tcpConn *TCPConn) Read(b []byte) (n int, err error) {
	return tcpConn.conn.Read(b)
}

func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

func (tcpConn *TCPConn) ReadMsg() (*common.WorldPacket, error) {
	headbuf := make([]byte, 4)
	// read len
	if _, err := io.ReadFull(tcpConn, headbuf); err != nil {
		return nil, err
	}
	tcpConn.Crypt.DecryptRecv(headbuf)
	msgLen := uint32(binary.LittleEndian.Uint16(headbuf[:2]))
	opCode := binary.LittleEndian.Uint16(headbuf[2:4])
	// check len
	if msgLen > tcpConn.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < HeaderLen {
		return nil, errors.New("message too short")
	}

	// data
	if msgLen > 4 {
		msgData := make([]byte, msgLen-4)
		if _, err := io.ReadFull(tcpConn, msgData); err != nil {
			return nil, err
		}
		packet := &common.WorldPacket{}
		packet.Initialize(opCode)
		packet.WriteBytes(msgData)
		return packet, nil
	}
	packet := &common.WorldPacket{}
	packet.Initialize(opCode)
	return packet, nil

}

func (tcpConn *TCPConn) WriteMsg(packet *common.WorldPacket) error {
	if tcpConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 5)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	tcpConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	tcpConn.Write(header.Bytes())
	return nil
}

func (tcpConn *TCPConn) InitCrypt(k []byte) {
	tcpConn.Crypt.Init(k)
}

func (tcpConn *TCPConn) IsTimeout(maxTime uint32) bool {
	var ret bool
	cur := uint32(time.Now().Unix() - tcpConn.activeTime)
	if cur > maxTime {
		ret = true
	}
	return ret
}
