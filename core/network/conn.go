package network

import (
	"net"

	"game_server/core/common"
)

type Conn interface {
	// Read() (int, []byte, error)
	Read(b []byte) (n int, err error)
	Write(b []byte)
	ReadMsg() (*common.WorldPacket, error)     //会对包头处理
	WriteMsg(packet *common.WorldPacket) error //会对包头处理
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	IsClosed() bool
	IsTimeout(maxTime uint32) bool
	InitCrypt(k []byte)
}
