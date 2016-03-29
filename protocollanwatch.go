package lanwatch

import (
	"github.com/huoyan108/gotcp"
	"log"
)

var (
	Illegal  uint16 = 0
	HalfPack uint16 = 255
	NoDo     uint16 = 254

	Login        uint16 = 1
	HeartBeat    uint16 = 2
	PosUp        uint16 = 53
	StationPosUp uint16 = 86
)

type LanWatchPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *LanWatchPacket) Serialize() []byte {
	switch this.Type {
	case Login:
		return this.Packet.(*LoginPacket).Serialize()
	case HeartBeat:
		return this.Packet.(*HeartPacket).Serialize()
	case NoDo:
		return this.Packet.(*NoDoPacket).Serialize()
	case PosUp:
		return this.Packet.(*PosUpPacket).Serialize()
	case StationPosUp:
		return this.Packet.(*MultiSationPacket).Serialize()
	}

	return nil
}

func NewLanWatchPacket(Type uint16, Packet gotcp.Packet) *LanWatchPacket {
	return &LanWatchPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type LanWatchProtocol struct {
}

func (this *LanWatchProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	log.Println("readPacket")
	smconn := c.GetExtraData().(*Conn)
	smconn.UpdateReadflag()

	buffer := smconn.GetBuffer()
	conn := c.GetRawConn()
	for {
		data := make([]byte, 2048)
		readLengh, err := conn.Read(data)
		log.Println(string(data))
		if err != nil {
			return nil, err
		}

		if readLengh == 0 {
			return nil, gotcp.ErrConnClosing
		} else {
			buffer.Write(data[0:readLengh])
			cmdid, pkglen := CheckProtocol(buffer)
			pkgbyte := make([]byte, pkglen)
			buffer.Read(pkgbyte)
			switch cmdid {
			case Login:
				pkg, daspkg := ParseLogin(pkgbyte, smconn)
				smconn.WriteToDas(daspkg)
				return NewLanWatchPacket(Login, pkg), nil
			case HeartBeat:
				pkg, daspkg := ParseHeart(pkgbyte)
				smconn.WriteToDas(daspkg)
				return NewLanWatchPacket(HeartBeat, pkg), nil
			case PosUp:
				pkg, daspkg := ParsePosUp(pkgbyte)
				smconn.WriteToDas(daspkg)
				return NewLanWatchPacket(PosUp, pkg), nil
			case StationPosUp:
				pkg, daspkg := ParseMultiStation(pkgbyte)
				smconn.WriteToDas(daspkg)
				return NewLanWatchPacket(StationPosUp, pkg), nil
			case NoDo:
				pkg := ParseNoDo(pkgbyte)
				return NewLanWatchPacket(NoDo, pkg), nil
			case Illegal:
			case HalfPack:
			}
		}
	}

}
