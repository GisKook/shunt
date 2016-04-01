package protocol

type DasHeartPacket struct {
	imei string
	batt string
}

func NewDasHeartPacket(imei string, batt string) *DasHeartPacket {
	return &DasHeartPacket{
		imei: imei,
		batt: batt,
	}
}

func (p *DasHeartPacket) Serialize() []byte {
	cmd := "$HSTAT:" + p.imei + "::" + p.batt + "\r\n"
	return []byte(cmd)
}
