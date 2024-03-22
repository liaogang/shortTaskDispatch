package shortid

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/jxskiss/base62"
	_ "github.com/jxskiss/base62"
	"hash/crc32"
	"testing"
	"time"
)

var table2 = [...]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

func weightConvert10to62New(input uint64) string {
	var m = make([]byte, 0)

	var table = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for {
		var t = input % 62
		input = input - t
		input /= 62

		var b = table[t]
		m = append(m, b)

		if input == 0 {
			break
		}

	}

	return string(m)
}

func weightConvert10to62NewNew() string {
	var m = [7]byte{}

	const length = 6
	var index = 0
	for index < length {

		//var b = table7[mRand.Intn(62)]
		var b = table7[1]

		m[index] = b

		index += 1
	}

	return string(m[:index])
}

func TestGen(t *testing.T) {

	var s = uuid.New()
	var str = s.String()

	var r = bytes.NewBuffer(s[:])
	var a uint64
	var b uint64
	binary.Read(r, binary.BigEndian, &a)
	binary.Read(r, binary.BigEndian, &b)

	println(weightConvert10to62New(a))
	println(weightConvert10to62New(b))

	//73WakrfVbNJBaAmhQtEeDv
	//kQn1G1CNJB67kiX1qNvwUd mine
	//3gsxpyfdv0kqn80y1p92bhakys

	_ = s
	print(str)
}

func Test_base64(t *testing.T) {

	var beg = time.Now()

	var i = 99999
	for i > 0 {
		i--

		var buf = make([]byte, 4)
		rand.Read(buf)
		base62.Encode(buf)
	}

	println("计时: ", time.Since(beg).String())

	var s = uuid.NewString()

	var c = crc32.ChecksumIEEE([]byte(s))

	beg = time.Now()
	i = 99999
	for i > 0 {
		i--

		weightConvert10to62(c)
	}

	println("计时: ", time.Since(beg).String())

	beg = time.Now()
	i = 99999
	for i > 0 {
		i--

		//weightConvert10to62NewNew()

		println(New()) //YzX1g1
	}

	println("计时: ", time.Since(beg).String())

	//println(string(ret))
	//print(string(ret))
}
