package aes

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/logxxx/utils"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext) == 0 {
		return []byte("")
	}

	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	if len(encrypt) == 0 {
		return []byte("")
	}

	srcLen := len(encrypt)
	paddingLen := int(encrypt[srcLen-1])
	if srcLen-paddingLen < 0 || paddingLen < 0 {
		return []byte("")
	}
	return encrypt[:srcLen-paddingLen]
}

func EncryptStr(key, src string) (encrypted string, err error) {
	resp, err := Encrypt([]byte(utils.MD5(key)), []byte(src))
	return string(resp), err
}

func Encrypt(key []byte, src []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	src = PKCS5Padding(src, block.BlockSize())
	encrypted = make([]byte, len(src))

	if len(encrypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks")
	}
	if len(encrypted) < len(src) {
		return nil, fmt.Errorf("crypto/cipher: output smaller than input")
	}

	ecb := NewECBEncrypter(block)
	ecb.CryptBlocks(encrypted, src)

	return
}

func Decrypt(key, src []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	decrypted = make([]byte, len(src))

	if len(decrypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks")
	}
	if len(decrypted) < len(src) {
		return nil, fmt.Errorf("crypto/cipher: output smaller than input")
	}

	ecb := NewECBDecrypter(block)
	ecb.CryptBlocks(decrypted, src)
	decrypted = PKCS5Trimming(decrypted)
	return
}
