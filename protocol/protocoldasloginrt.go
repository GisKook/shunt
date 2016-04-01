package protocol

import (
	"bytes"
)

type LoginRt struct {
	Result bool
}

func ParseDasLoginRt(buffer []byte) *LoginRt {
	flag := []byte{':'}
	res := bytes.Split(buffer, flag)
	if res[3][0] == '1' {
		return &LoginRt{
			Result: true,
		}
	}

	return &LoginRt{
		Result: false,
	}
}
