package lanwatch

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type DasPosupPacket struct {
	Uid   string
	Timen uint64
	Batt  string
	Speed float64
	Long  float64
	Lat   float64
}

func NewDasPosupPacket(uid string, timen uint64, batt string, speed float64, long float64, lat float64) *DasPosupPacket {
	return &DasPosupPacket{
		Uid:   uid,
		Timen: timen,
		Batt:  batt,
		Speed: speed,
		Long:  long,
		Lat:   lat,
	}
}

func (p *DasPosupPacket) Serialize() []byte {
	speed := strconv.FormatFloat(p.Speed, 'f', 1, 32)
	long := strconv.FormatFloat(p.Long, 'f', 6, 32)
	lat := strconv.FormatFloat(p.Lat, 'f', 6, 32)

	tm := time.Unix(int64(p.Timen), 0)
	log.Println(tm)
	gpstime := fmt.Sprintf(":%02d%02d%02d-%02d%02d%02d:", tm.Year()-2000, tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	log.Println(gpstime)
	cmd := "$POSUP:" + p.Uid + gpstime + p.Batt + ":" + speed + ":0:2:0:" + long + "," + lat + "\r\n"

	log.Println(cmd)

	return []byte(cmd)
}
