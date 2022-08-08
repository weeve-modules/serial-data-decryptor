package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Data struct {
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
	IV        string `json:"iv"`
}

type PlainFormatter struct {
	TimestampFormat string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	return []byte(fmt.Sprintf("%s %s : %s\n", timestamp, entry.Level, entry.Message)), nil
}

func GetEnvAsserted(envVarName string) string {
	var thisEnvVar = os.Getenv(envVarName)
	if len(thisEnvVar) == 0 {
		log.Fatal(envVarName, " was not found in the current environment")
	}
	return thisEnvVar
}

func Decrypt(iv, ct []byte) ([]byte, error) {
	var key = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, iv, ct, nil)
}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	msg := Data{}
	dec.Decode(&msg)
	fmt.Printf("Received : %+v", msg)
	iv, err := base64.RawStdEncoding.DecodeString(msg.IV)
	if err != nil {
		log.Error(err)
	}
	ct, err := base64.RawStdEncoding.DecodeString(msg.Data)
	if err != nil {
		log.Error(err)
	}
	result, err := Decrypt(iv, ct)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("Decrypted : %+v", result)

	conn.Close()
}

func main() {
	// module_name := GetEnvAsserted("MODULE_NAME")
	// ingress_host := GetEnvAsserted("INGRESS_HOST")
	// ingress_port := GetEnvAsserted("INGRESS_PORT")
	// egess_urls := GetEnvAsserted("EGRESS_URLS")
	// input_label := GetEnvAsserted("INPUT_LABEL")

	// log.Info("%s running on %s at port %s with end-point set to %s", module_name, ingress_host, ingress_port, egess_urls)
	// log.Info("keys to decrypt: %s", input_label)

	fmt.Println("start the server")
	ln, err := net.Listen("tcp", ":"+"80")
	if err != nil {
		log.Error(err)
	}

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			log.Error(err)
			continue
		}
		go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}
