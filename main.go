package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Data map[string]interface{}

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

func sendData(data []byte) {
	egress_Urls := strings.Split(GetEnvAsserted("EGRESS_URLS"), ",")
	for _, egress_Url := range egress_Urls {
		fmt.Println("start client")
		// Create a HTTP post request
		r, err := http.NewRequest("POST", egress_Url, bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err)
			continue
		}

		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer res.Body.Close()

		post := &Data{}
		derr := json.NewDecoder(res.Body).Decode(post)
		if derr != nil {
			fmt.Println(err)
		}

		if res.StatusCode != http.StatusOK {
			fmt.Println(err)
			continue
		}
		fmt.Println("sent")
	}
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

	result, err := aesgcm.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	msg := Data{}
	dec.Decode(&msg)
	fmt.Printf("Received : %+v\n\n", msg)

	iv, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["data"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	ct, err := base64.StdEncoding.DecodeString(fmt.Sprint(msg["iv"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	d, err := Decrypt(iv, ct)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("Decrypted value")
	fmt.Printf("+ %s\n", d)
	conn.Close()

	sendData(d)
}

func startServer() {
	fmt.Println("start the server")
	ln, err := net.Listen("tcp", GetEnvAsserted("INGRESS_HOST")+":"+GetEnvAsserted("INGRESS_PORT"))
	fmt.Println("started the server: " + ln.Addr().String())
	if err != nil {
		log.Fatal(err)
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

func main() {
	godotenv.Load("./docker/docker.env") // TODO: Only use it for testing locally
	log.Info("%s running on %s at port %s with end-point set to %s for data %s", GetEnvAsserted("MODULE_NAME"), GetEnvAsserted("INGRESS_HOST"), GetEnvAsserted("INGRESS_PORT"), GetEnvAsserted("EGRESS_URLS"), GetEnvAsserted("INPUT_LABELS"))

	startServer()
}
