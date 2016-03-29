package lanwatch

import (
	"log"
)

type DasLoginPacket struct {
	Uid string
}

func NewDasLoginPacket(uid string) *DasLoginPacket {
	return &DasLoginPacket{
		Uid: uid,
	}
}

func (p *DasLoginPacket) Serialize() []byte {
	cmd := "$LOGIN:" + p.Uid + ":DK-PE100:V2.0-20150927\r\n"
	log.Println(cmd)
	return []byte(cmd)
}
