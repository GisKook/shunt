package lanwatch

import (
	"bytes"
	//	"encoding/binary"
	//	"errors"
	"log"
	"strconv"
)

func CheckSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	log.Println("cp 1")
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	log.Println("cp 2")
	//log.Println(string(buffer.Bytes()))
	//if buffer.Bytes()[0] != '[' {
	//	return Illegal, 0
	//}
	//find [
	leftFlag := []byte{'['}
	l := bytes.Index(buffer.Bytes(), leftFlag)
	if l == -1 {
		return Illegal, 0
	}
	log.Println("cp 3")
	p := 0
	bFlag := false
	for {
		if p < bufferlen && buffer.Bytes()[p] != ']' {
			p++
		} else {
			bFlag = true
			break
		}
	}
	if !bFlag {
		return Illegal, 0
	}
	p++
	if l > p {
		return Illegal, 0
	}
	var source []byte
	source = buffer.Bytes()[0:p]
	flag := []byte{','}
	res := bytes.Split(source, flag)

	flag = []byte{']'}
	res1 := bytes.Split(res[8], flag)
	protocType := string(res1[0][1:])

	nType, _ := strconv.ParseInt(protocType, 10, 32)

	log.Println("type", nType)
	switch uint16(nType) {
	case Login:
		return Login, uint16(p)
	case HeartBeat:
		return HeartBeat, uint16(p)
	case PosUp:
		return PosUp, uint16(p)
	case StationPosUp:
		return StationPosUp, uint16(p)
	default:
		return NoDo, uint16(p)
	}
	//var source []byte
	//source = buffer.Bytes()[0:p]
	//t1flag := []byte{'T', '1'}
	//t2flag := []byte{'T', '2'}
	//t53flag := []byte{'T', '5', '3'}
	//if bytes.Contains(source, t1flag) {
	//	return Login, uint16(p)
	//} else if bytes.Contains(source, t2flag) {
	//	return HeartBeat, uint16(p)
	//} else if bytes.Contains(source, t53flag) {
	//	return PosUp, uint16(p)
	//} else {
	//	return NoDo, uint16(p)
	//}

	return Illegal, 0
}

//func GetGatewayID(buffer []byte) (uint64, *bytes.Reader) {
//	reader := bytes.NewReader(buffer)
//	reader.Seek(5, 0)
//	uid := make([]byte, 6)
//	reader.Read(uid)
//	gid := []byte{0, 0}
//	gid = append(gid, uid...)
//	return binary.BigEndian.Uint64(gid), reader
//}
//
//func CheckNsqProtocol(message []byte) (uint64, uint32, *Report.Command, error) {
//	command := &Report.ControlReport{}
//	err := proto.Unmarshal(message, command)
//	if err != nil {
//		log.Println("unmarshal error")
//		return 0, 0, nil, errors.New("unmarshal error")
//	} else {
//		gatewayid := command.Tid
//		serialnum := command.SerialNumber
//		cmd := command.GetCommand()
//
//		return gatewayid, serialnum, cmd, nil
//	}
//}
