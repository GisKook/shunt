package lanwatch

import (
	"bytes"
	"log"
	"strconv"
	"time"
)

type StationPacket struct {
	Lac    uint32
	Cellid uint32
	Rxlevx uint32
}
type MultiSationPacket struct {
	Uid     uint64
	Head    string
	Voltage string
	Num     uint16
	Timen   uint32
}

func (p *MultiSationPacket) Serialize() []byte {
	var sBuf string
	sBuf = "[V1.0.0,"
	sBuf += p.Head
	sBuf += ",1,abcd,"
	sBuf += time.Now().Format("2006-01-02 15:04:05")
	sBuf += ","
	sBuf += strconv.Itoa(int(p.Uid))
	sBuf += ",S86]"
	var buf []byte
	buf = []byte(sBuf)

	log.Println("s86", sBuf)
	return buf
}

func NewMultiSationPacket(Uid uint64, Head string, Voltage string, Num uint16, Timen uint32) *MultiSationPacket {
	return &MultiSationPacket{
		Uid:     Uid,
		Head:    Head,
		Voltage: Voltage,
		Num:     Num,
		Timen:   Timen,
	}
}

func ParseMultiStation(buffer []byte) (*MultiSationPacket, *DasStationPosupPacket) {
	flag := []byte{','}
	res := bytes.Split(buffer, flag)
	gatewayid, _ := strconv.ParseInt(string(res[6]), 10, 64)
	head := string(res[1])
	nVoltage, _ := strconv.ParseInt(string(res[9]), 10, 32)
	voltage := strconv.Itoa(int(nVoltage) * 2 * 10)

	flag = []byte{']'}
	cTimen := bytes.Split(res[11], flag)
	timen, _ := strconv.ParseInt(string(cTimen[0]), 10, 64)
	log.Println("timen", timen)
	flag[0] = '#'
	stationsSource := bytes.Split(res[10], flag)
	//log.Println(stationsSource)

	flag[0] = '$'
	source := bytes.Split(stationsSource[0], flag)
	num, _ := strconv.ParseInt(string(source[1]), 10, 32)
	//log.Println("num", num)

	var sStations string
	flag[0] = '|'
	if num > 4 {
		num = 4
	}
	var temp int64
	for i := 0; i < int(num); i++ {
		stationInfo := bytes.Split(stationsSource[i+1], flag)
		if i == 0 {
			temp, _ = strconv.ParseInt(string(stationInfo[0]), 16, 32)
			sStations += strconv.Itoa(int(temp))
			sStations += ","
			temp, _ = strconv.ParseInt(string(stationInfo[1]), 16, 32)
			sStations += strconv.Itoa(int(temp))
		} else {
			sStations += ":"
			temp, _ = strconv.ParseInt(string(stationInfo[0]), 16, 32)
			sStations += strconv.Itoa(int(temp))
			sStations += ","
			temp, _ = strconv.ParseInt(string(stationInfo[1]), 16, 32)
			sStations += strconv.Itoa(int(temp))
			sStations += ","
			temp, _ = strconv.ParseInt(string(stationInfo[3]), 10, 32)
			sStations += strconv.Itoa(int(temp))
		}
	}
	if num < 4 && num > 0 {
		for i := 0; i < 4-int(num); i++ {
			sStations += ":,,"
		}
	}
	//log.Println(sStations)
	return NewMultiSationPacket(uint64(gatewayid), head, voltage, uint16(num), uint32(timen)), NewDasStationPosupPacket(string(res[6]), uint64(timen), voltage, 0, sStations)
}
