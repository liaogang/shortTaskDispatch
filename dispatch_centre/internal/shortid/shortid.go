package shortid

import (
	"github.com/google/uuid"
	"hash/crc32"
)

func New() string {
	var s = uuid.NewString()

	var c = crc32.ChecksumIEEE([]byte(s))

	return weightConvert10to62(c)
}

var table7 = [...]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

func weightConvert10to62(input uint32) string {
	var m = [20]byte{}

	var index = 0
	for {
		var t = input % 62
		input = input - t
		input /= 62

		var b = table7[t]

		m[index] = b

		index += 1

		if input == 0 {
			break
		}

	}

	return string(m[:index])
}
