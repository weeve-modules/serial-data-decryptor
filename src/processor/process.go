package processor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"serial-data-decryptor/models"
	"serial-data-decryptor/utility"

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

	if strings.TrimSpace(msg.IV) == "" {
		errorList = append(errorList, "iv is empty/nil")
	}

	if strings.TrimSpace(msg.Cyphertext) == "" {
		errorList = append(errorList, "cyphertext is empty/nil")
	}

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "; "))
	} else {
		return nil
	}
}

func decrypt(msg models.Data) ([]byte, error) {
	log.Debug("Decrypting...")

	key, err := base64.StdEncoding.DecodeString(utility.GetEnvAsserted("AES_KEY"))
	if err != nil {
		return nil, err
	}
	log.Debug("key = ", string(key))

	iv, err := base64.StdEncoding.DecodeString(msg.IV)
	if err != nil {
		return nil, err
	}
	log.Debug("iv = ", string(iv))

	ct, err := base64.StdEncoding.DecodeString(msg.Cyphertext)
	if err != nil {
		return nil, err
	}
	log.Debug("cypthertext = ", string(ct))

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
