package protocol

import (
	"bytes"
)

type HeartPacket struct {
	manufacturer string
	imei         string
	batt         string
}

func (p *HeartPacket) Serialize() []byte {
	var feedback string
	feedback = p.manufacturer + "*" + p.imei + "*0002*LK]"
	buf := []byte(feedback)

	return buf
}

func NewHeartPacket(manufacturer string, imei string, batt string) *HeartPacket {
	return &HeartPacket{
		manufacturer: manufacturer,
		imei:         imei,
		batt:         batt,
	}
}

func ParseHeart(buffer []byte) (*HeartPacket, *DasHeartPacket, string) {
	flag := []byte{'*'}
	res := bytes.Split(buffer, flag)
	manufacturer := string(res[0])
	imei := string(res[1])

	commaflag := []byte{','}
	batts := bytes.Split(buffer, commaflag)
	var batt string
	if len(batts) == 4 {
		batt = string(batts[3])
	}

	return NewHeartPacket(manufacturer, imei, batt), NewDasHeartPacket(imei, batt), batt
}
