package lanwatch

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

type HeartPacket struct {
	Uid  uint64
	Head string
}

func (this *HeartPacket) Serialize() []byte {
	var sBuf string
	sBuf = "[V1.0.0,"
	sBuf += this.Head
	sBuf += ",1,abcd,"
	sBuf += time.Now().Format("2006-01-02 15:04:05")
	sBuf += ","
	sBuf += strconv.Itoa(int(this.Uid))
	sBuf += ",S2]"
	var buf []byte
	buf = []byte(sBuf)
	log.Println(string(buf))
	return buf
}

func NewHeartPacket(Uid uint64, Head string) *HeartPacket {
	return &HeartPacket{
		Uid:  Uid,
		Head: Head,
	}
}

func ParseHeart(buffer []byte) (*HeartPacket, *DasHeartPacket) {
	flag := []byte{','}
	res := bytes.Split(buffer, flag)
	id := res[6]

	gatewayid, _ := strconv.ParseInt(string(id), 10, 64)
	return NewHeartPacket(uint64(gatewayid), string(res[1])), NewDasHeartPacket(string(id), "0", "60")
}
