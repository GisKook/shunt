package lanwatch

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

type NoDoPacket struct {
	Uid        uint64
	Head       string
	ProtocType string
}

func (p *NoDoPacket) Serialize() []byte {
	var sBuf string
	sBuf = "[V1.0.0,"
	sBuf += p.Head
	sBuf += ",1,abcd,"
	sBuf += time.Now().Format("2006-01-02 15:04:05")
	sBuf += ","
	sBuf += strconv.Itoa(int(p.Uid))
	sBuf += ",S"
	sBuf += p.ProtocType
	sBuf += "]"
	var buf []byte
	buf = []byte(sBuf)
	log.Println(sBuf)
	return buf
}

func NewNoDoPakcet(Uid uint64, Head string, ProtocType string) *NoDoPacket {
	return &NoDoPacket{
		Uid:        Uid,
		Head:       Head,
		ProtocType: ProtocType,
	}
}

func ParseNoDo(buffer []byte) *NoDoPacket {
	flag := []byte{','}
	res := bytes.Split(buffer, flag)
	id := res[6]
	gatewayid, _ := strconv.ParseInt(string(id), 10, 64)

	protocType := string(res[8][1:])
	return NewNoDoPakcet(uint64(gatewayid), string(res[1]), protocType)

}
