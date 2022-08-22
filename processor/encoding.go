package processor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"serial-data-decryptor/models"
)

func Process(msg models.Data) ([]byte, error) {
	return Decrypt(msg)
}

func Decrypt(msg models.Data) ([]byte, error) {
	var key = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	iv, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["data"]))
	if err != nil {
		return nil, err
	}
	ct, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["iv"]))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	result, err := aesgcm.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Decrypted message: %s\n", result)

	return result, nil
}
