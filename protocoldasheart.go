package lanwatch

import (
	"log"
)

type DasHeartPacket struct {
	Uid    string
	Status string
	Batt   string
}

func NewDasHeartPacket(uid string, status string, batt string) *DasHeartPacket {
	return &DasHeartPacket{
		Uid:    uid,
		Status: status,
		Batt:   batt,
	}
}

func (p *DasHeartPacket) Serialize() []byte {
	cmd := "$HSTAT:" + p.Uid + ":" + p.Status + ":" + p.Batt + "\r\n"
	log.Println(cmd)

	return []byte(cmd)
}
