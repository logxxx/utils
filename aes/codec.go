package aes

import (
	"bytes"
	"encoding/gob"
)

type AesCodec struct {
	EncryptKey string
}

func (c AesCodec) Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return Encrypt([]byte(c.EncryptKey), b.Bytes())
}

func (c AesCodec) Unmarshal(b []byte, v interface{}) error {
	b, err := Decrypt([]byte(c.EncryptKey), b)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}

func (c AesCodec) Name() string {
	return "aes"
}
