package shortid

import (
	"bytes"
	"crypto/rand"
)

var table = [...]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

/*
生成5个字符的id, 生产环境先检查重复性
*/

func New() string {

	var buf = make([]byte, 5)
	rand.Read(buf)

	var collect = bytes.NewBuffer(nil)

	for _, b := range buf {
		c := b % 62
		collect.WriteByte(table[c])
	}

	return string(collect.Bytes())
}
