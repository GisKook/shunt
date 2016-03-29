package lanwatch

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

type PosUpPacket struct {
	Uid       uint64
	DataStyle string
	Long      float32
	Lat       float32
	Height    float32
	Speed     float32
	Direction float32
	Timen     uint64
	Head      string
}

func (p *PosUpPacket) Serialize() []byte {
	var sBuf string
	sBuf = "[V1.0.0,"
	sBuf += p.Head
	sBuf += ",1,abcd,"
	sBuf += time.Now().Format("2006-01-02 15:04:05")
	sBuf += ","
	sBuf += strconv.Itoa(int(p.Uid))
	sBuf += ",S2]"
	var buf []byte
	buf = []byte(sBuf)

	return buf
}

func NewPosUpPacket(Uid uint64, DataStyle string, Long float32, Lat float32, Height float32, Speed float32, Direction float32, Timen uint64, Head string) *PosUpPacket {
	return &PosUpPacket{
		Uid:       Uid,
		DataStyle: DataStyle,
		Long:      Long,
		Lat:       Lat,
		Height:    Height,
		Speed:     Speed,
		Direction: Direction,
		Timen:     Timen,
		Head:      Head,
	}
}

func ParsePosUp(buffer []byte) (*PosUpPacket, *DasPosupPacket) {
	flag := []byte{','}
	res := bytes.Split(buffer, flag)
	id := res[6]
	gatewayid, _ := strconv.ParseInt(string(id), 10, 64)
	dataStyle := string(res[9])
	long, _ := strconv.ParseFloat(string(res[10]), 32)
	lat, _ := strconv.ParseFloat(string(res[11]), 32)
	height, _ := strconv.ParseFloat(string(res[12]), 32)
	speed, _ := strconv.ParseFloat(string(res[13]), 32)
	direction, _ := strconv.ParseFloat(string(res[14]), 32)

	flag = []byte{']'}
	cTimen := bytes.Split(res[15], flag)
	log.Println(cTimen[0])
	timen, _ := strconv.ParseInt(string(cTimen[0]), 10, 64)
	//timen, _ := strconv.ParseInt(string(res[15]), 10, 64)

	log.Println(gatewayid, dataStyle, long, lat, height, speed, direction, timen)
	return NewPosUpPacket(uint64(gatewayid), dataStyle, float32(long), float32(lat), float32(height), float32(speed), float32(direction), uint64(timen), string(res[1])), NewDasPosupPacket(string(id), uint64(timen), "60", speed, long, lat)
}
