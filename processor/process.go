package processor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"serial-data-decryptor/models"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Process(msg models.Data) ([]byte, error) {
	err := validateData(msg)
	if err != nil {
		return nil, err
	}

	return decrypt(msg)
}

func validateData(msg models.Data) error {
	var errorList []string

	iv := msg["iv"]
	if iv == nil || (strings.TrimSpace(iv.(string)) == "") {
		errorList = append(errorList, "iv is empty/nil")
	}

	data := msg["cyphertext"]
	if data == nil || (strings.TrimSpace(data.(string)) == "") {
		errorList = append(errorList, "cyphertext is empty/nil")
	}

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList[:], ","))
	} else {
		return nil
	}
}

func decrypt(msg models.Data) ([]byte, error) {
	var key = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	iv, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["iv"]))
	if err != nil {
		return nil, err
	}
	ct, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["cyphertext"]))
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

	if len(iv) != aesgcm.NonceSize() {
		return nil, fmt.Errorf("crypto/cipher: incorrect nonce length given to GCM")
	}

	result, err := aesgcm.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, err
	}

	log.Infof("Decrypted message: %s\n", result)

	return result, nil
}
