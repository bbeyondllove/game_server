package network

import (
	"game_server/core/logger"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	sync.Mutex
	Addr            string
	ConnNum         int
	ConnectInterval time.Duration
	AutoReconnect   bool
	NewAgent        func(*TCPConn) Agent
	conns           ConnSet
	wg              sync.WaitGroup
	closeFlag       bool

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
}

func (client *TCPClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		logger.Infof("invalid ConnNum, reset to ", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		logger.Infof("invalid ConnectInterval, reset to ", client.ConnectInterval)
	}
	if client.NewAgent == nil {
		logger.Fatal("NewAgent must not be nil")
	}
	if client.conns != nil {
		logger.Fatal("client is running")
	}

	client.conns = make(ConnSet)
	client.closeFlag = false
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag {
			return conn
		}

		logger.Infof("connect to  error:", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}
	client.conns[conn] = struct{}{}
	client.Unlock()

	tcpConn := newTCPConn(conn, client.MaxMsgLen)
	agent := client.NewAgent(tcpConn)
	agent.Run()

	// cleanup
	tcpConn.Close()
	client.Lock()
	delete(client.conns, conn)
	client.Unlock()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *TCPClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	client.wg.Wait()
}

func (client *TCPClient) IsClosed() bool {
	return client.closeFlag
}
