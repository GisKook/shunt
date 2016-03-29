package lanwatch

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type DasStationPosupPacket struct {
	Uid      string
	Timen    uint64
	Batt     string
	Speed    float64
	Stations string
}

func NewDasStationPosupPacket(uid string, timen uint64, batt string, speed float64, Stations string) *DasStationPosupPacket {
	return &DasStationPosupPacket{
		Uid:      uid,
		Timen:    timen,
		Batt:     batt,
		Speed:    speed,
		Stations: Stations,
	}
}

func (p *DasStationPosupPacket) Serialize() []byte {
	speed := strconv.FormatFloat(p.Speed, 'f', 1, 32)

	tm := time.Unix(int64(p.Timen), 0)
	log.Println(tm)
	gpstime := fmt.Sprintf(":%02d%02d%02d-%02d%02d%02d:", tm.Year()-2000, tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	log.Println(gpstime)
	cmd := "$POSUP:" + p.Uid + gpstime + p.Batt + ":" + speed + ":0:2:1:" + p.Stations + "\r\n"

	log.Println(cmd)

	return []byte(cmd)
}
