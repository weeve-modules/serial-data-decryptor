package processor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"serial-data-decryptor/models"
	"serial-data-decryptor/utility"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)
}

func Process(msg models.Data) ([]byte, error) {
	err := validate.Struct(msg)
	if err != nil {
		return nil, err
	}

	return decrypt(msg)
}

func decrypt(msg models.Data) ([]byte, error) {
	log.Debug("Decrypting...")

	key, err := base64.StdEncoding.DecodeString(utility.GetEnvAsserted("AES_KEY"))
	if err != nil {
		return nil, err
	}

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
