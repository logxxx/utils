package randutil

import (
	"encoding/hex"
	"math/rand"
	"time"
)

var letters = []byte("abcdefghjkmnpqrstuvwxyz123456789")
var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // 61

func init() {
	rand.Seed(time.Now().Unix())
}

// RandLow 随机字符串，包含 1~9 和 a~z - [i,l,o]
func RandLow(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	for i, x := range b {
		arc = x & 31
		b[i] = letters[arc]
	}
	return string(b)
}

// RandUp 随机字符串，包含 英文字母和数字
func RandStr(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	for i, x := range b {
		arc = x & 61
		b[i] = longLetters[arc]
	}
	return string(b)
}

// RandHex 生成16进制格式的随机字符串
func RandHex(n int) string {
	if n <= 0 {
		return ""
	}
	var need int
	if n&1 == 0 { // even
		need = n
	} else { // odd
		need = n + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return ""
	}
	hex.Encode(dst, src)
	return string(dst[:n])
}
