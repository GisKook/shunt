package protocol

import (
	"bytes"
)

type DasPosUpPacket struct {
	IMEI      string
	Time      string
	Batt      string
	Speed     string
	Parse     string
	PosReason string
	Postype   string

	Longitude string
	Latitude  string
	Lac       string
	Cid       string
	Lac1      string
	Cid1      string
	Dbm1      string
	Lac2      string
	Cid2      string
	Dbm2      string
	Lac3      string
	Cid3      string
	Dbm3      string
}

func (p *DasPosUpPacket) Serialize() []byte {
	cmd := "$POSUP:" + p.IMEI + ":" + p.Time + ":" + p.Batt + ":" + p.Speed + ":" + p.Parse + ":" + p.PosReason + ":" + p.Postype + ":"
	if p.Postype == "0" {
		cmd += p.Longitude + "," + p.Latitude + "\r\n"
	} else {
		cmd += p.Lac + "," + p.Cid + ":" + p.Lac1 + "," + p.Cid1 + "," + p.Dbm1 + ":"
		cmd += p.Lac2 + "," + p.Cid2 + "," + p.Dbm2 + ":"
		cmd += p.Lac3 + "," + p.Cid3 + "," + p.Dbm3 + "\r\n"
	}

	return []byte(cmd)
}

func ParsePosUp(buffer []byte) *DasPosUpPacket {
	flag := []byte{'*'}
	res_imei := bytes.Split(buffer, flag)
	imei := string(res_imei[1])

	flag = []byte{','}
	res := bytes.Split(buffer, flag)
	time := string(res[1][4:6]) + string(res[1][2:4]) + string(res[1][0:2])
	time += "-" + string(res[2][0:2]) + string(res[2][2:4]) + string(res[2][4:6])
	batt := string(res[13])
	speed := string(res[8])
	parse := "0"
	posreason := "2"
	var postype string = "0"
	if res[3][0] == 'V' {
		postype = "1"
	}
	latitude := string(res[4])
	longitude := string(res[6])
	lac := string(res[21])
	cid := string(res[22])
	lac1 := string(res[24])
	cid1 := string(res[25])
	dbm1 := string(res[26])
	lac2 := string(res[27])
	cid2 := string(res[28])
	dbm2 := string(res[29])
	lac3 := string(res[30])
	cid3 := string(res[31])
	dbm3 := string(res[32])

	return &DasPosUpPacket{
		IMEI:      imei,
		Time:      time,
		Batt:      batt,
		Speed:     speed,
		Parse:     parse,
		PosReason: posreason,
		Postype:   postype,
		Longitude: longitude,
		Latitude:  latitude,
		Lac:       lac,
		Cid:       cid,
		Lac1:      lac1,
		Cid1:      cid1,
		Dbm1:      dbm1,
		Lac2:      lac2,
		Cid2:      cid2,
		Dbm2:      dbm2,
		Lac3:      lac3,
		Cid3:      cid3,
		Dbm3:      dbm3,
	}
}
