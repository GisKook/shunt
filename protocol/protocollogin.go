package protocol

import (
	"bytes"
	//	"encoding/binary"
	"strconv"
)

type LoginPacket struct {
	manufacturer string
	imei         string
}

func (p *LoginPacket) Serialize() []byte {
	var feedback string
	feedback = p.manufacturer + "*" + p.imei + "*0002*LK]"
	buf := []byte(feedback)

	return buf
}

func NewLoginPakcet(manufacturer string, imei string) *LoginPacket {
	return &LoginPacket{
		manufacturer: manufacturer,
		imei:         imei,
	}
}

func ParseLogin(buffer []byte) (*LoginPacket, *DasLoginPacket, uint64, string) {
	flag := []byte{'*'}
	res := bytes.Split(buffer, flag)
	manufacturer := string(res[0])
	imei := string(res[1])
	trackerid, _ := strconv.ParseUint(string(imei), 10, 64)

	commaflag := []byte{','}
	batt := bytes.Split(res[3], commaflag)
	var batt_value string = "0"
	if len(batt) == 4 {
		batt_value = string(batt[3])
	}

	return NewLoginPakcet(manufacturer, imei), NewDasLoginPacket(imei), trackerid, batt_value
}
