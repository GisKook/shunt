package lanwatch

import (
	"bytes"
	"github.com/huoyan108/gotcp"
	"log"
	"time"

	"net"
)

var ConnSuccess uint8 = 0
var ConnUnauth uint8 = 1

type ConnConfig struct {
	HeartBeat    uint8
	ReadLimit    int64
	WriteLimit   int64
	NsqChanLimit int32
}

type Conn struct {
	conn          *gotcp.Conn
	config        *ConnConfig
	recieveBuffer *bytes.Buffer
	ticker        *time.Ticker
	readflag      int64
	writeflag     int64
	closeChan     chan bool
	index         uint32
	uid           uint64
	status        uint8

	dasconn    *net.TCPConn
	dasCmdChan chan gotcp.Packet
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", "192.168.2.224:1725")
	dasconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Printf("conn to das fail %s\n", err.Error())
	}
	return &Conn{
		conn:          conn,
		recieveBuffer: bytes.NewBuffer([]byte{}),
		config:        config,
		readflag:      time.Now().Unix(),
		writeflag:     time.Now().Unix(),
		ticker:        time.NewTicker(time.Duration(config.HeartBeat) * time.Second),
		closeChan:     make(chan bool),
		index:         0,
		status:        ConnUnauth,
		dasconn:       dasconn,
		dasCmdChan:    make(chan gotcp.Packet, 64),
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.recieveBuffer.Reset()
	close(c.closeChan)
	close(c.dasCmdChan)
}

func (c *Conn) GetGatewayID() uint64 {
	return c.uid
}
func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

func (c *Conn) sendToDas() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case p := <-c.dasCmdChan:
			if p != nil {
				c.dasconn.Write(p.Serialize())
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) WriteToDas(p gotcp.Packet) {
	c.dasCmdChan <- p
}

func (c *Conn) UpdateReadflag() {
	c.readflag = time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.writeflag = time.Now().Unix()
}

func (c *Conn) SetStatus(status uint8) {
	c.status = status
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	var now int64
	for {
		select {
		case <-c.ticker.C:
			now = time.Now().Unix()
			if now-c.readflag > c.config.ReadLimit {
				log.Println("read linmit")
				return
			}
			if now-c.writeflag > c.config.WriteLimit {
				log.Println("write limit")
				return
			}
			if c.status == ConnUnauth {
				log.Printf("unauth's gateway gatewayid %d\n", c.uid)
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) recvdas() {
	for {
		buffer := make([]byte, 1024)
		c.dasconn.Read(buffer)
		log.Println(string(buffer))
	}
}

func (c *Conn) Do() {
	//go c.checkHeart()
	go c.recvdas()
	go c.sendToDas()
	//go c.writeToclientLoop()
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	log.Println("new conn")
	heartbeat := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		HeartBeat:  uint8(heartbeat),
		ReadLimit:  int64(readlimit),
		WriteLimit: int64(writelimit),
	}
	//log.Println(heartbeat,readlimit,writelimit)
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	log.Println("closeconn")
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn.GetGatewayID())
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	lanpacket := p.(*LanWatchPacket)
	switch lanpacket.Type {
	case Login:
		c.AsyncWritePacket(lanpacket, time.Second)
	case HeartBeat:
		c.AsyncWritePacket(lanpacket, time.Second)
	case PosUp:
		c.AsyncWritePacket(lanpacket, time.Second)
	case StationPosUp:
		c.AsyncWritePacket(lanpacket, time.Second)
	case NoDo:
		c.AsyncWritePacket(lanpacket, time.Second)
	}

	return true
}
