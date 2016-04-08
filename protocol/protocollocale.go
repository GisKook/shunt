package protocol

import (
	"fmt"
)

type LocalePacket struct {
	manufacturer string
	imei         uint64
}

func (p *LocalePacket) Serialize() []byte {
	locale := "[" + p.manufacturer + "*" + fmt.Sprint(p.imei) + "*LZ,1,8]"

	return []byte(locale)
}

func NewLocalePacket(manufacturer string, imei uint64) *LocalePacket {
	return &LocalePacket{
		manufacturer: manufacturer,
		imei:         imei,
	}
}
