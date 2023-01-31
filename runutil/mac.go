package runutil

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMacAddrsMd5(salt string) (macAddrsMd5 []string) {
	macAddrs := GetMacAddrs()
	for _, v := range macAddrs {
		macAddrsMd5 = append(macAddrsMd5, md5v(v+salt))
	}
	return macAddrsMd5
}

func md5v(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
