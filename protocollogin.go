package lanwatch

import (
	"bytes"
	//	"encoding/binary"
	"log"
	"strconv"
	"time"
)

var (
	Infrared      uint8 = 0
	DoorMagnetic  uint8 = 1
	WarningButton uint8 = 2
)

type Device struct {
	Oid     uint64
	Type    uint8
	Company uint16
	Status  uint8
	Name    string
}

type LoginPacket struct {
	Uid    uint64
	Head   string
	Result uint8
}

func (p *LoginPacket) Serialize() []byte {
	var sBuf string
	sBuf = "[V1.0.0,"
	sBuf += p.Head
	sBuf += ",1,abcd,"
	sBuf += time.Now().Format("2006-01-02 15:04:05")
	sBuf += ","
	sBuf += strconv.Itoa(int(p.Uid))
	sBuf += ",S1]"
	var buf []byte
	buf = []byte(sBuf)
	log.Println(sBuf)
	return buf
}

func NewLoginPakcet(Uid uint64, Head string, Result uint8) *LoginPacket {
	return &LoginPacket{
		Uid:    Uid,
		Head:   Head,
		Result: Result,
	}
}

func ParseLogin(buffer []byte, c *Conn) (*LoginPacket, *DasLoginPacket) {
	flag := []byte{','}
	res := bytes.Split(buffer, flag)
	id := res[6]
	gatewayid, _ := strconv.ParseInt(string(id), 10, 64)
	c.uid = uint64(gatewayid)
	log.Println("uid", gatewayid)
	c.SetStatus(ConnSuccess)
	NewConns().Add(c)
	log.Println("addnewconn")
	return NewLoginPakcet(uint64(gatewayid), string(res[1]), 1), NewDasLoginPacket(string(id))

}
