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

var inputLabels []string

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

	result, err := aesgcm.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handleConnection(conn net.Conn, egressUrl string) {
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

	sendData(d, egressUrl)
}

func sendData(data []byte, egressUrl string) {
	fmt.Println("start client")
	// Create a HTTP post request
	r, err := http.NewRequest("POST", egressUrl, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	post := &Data{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != http.StatusCreated {
		fmt.Println(err)
		return
	}
	fmt.Println("sent status: " + r.Response.Status)
}

func main() {
	godotenv.Load("./docker/docker.env") // TODO: Only use it for testing locally
	module_name := GetEnvAsserted("MODULE_NAME")
	ingress_host := GetEnvAsserted("INGRESS_HOST")
	ingress_port := GetEnvAsserted("INGRESS_PORT")
	egess_urls := GetEnvAsserted("EGRESS_URLS")
	input_labels := GetEnvAsserted("INPUT_LABELS")

	log.Info("keys to decrypt: %s", input_labels)
	log.Info("%s running on %s at port %s with end-point set to %s", module_name, ingress_host, ingress_port, egess_urls)

	inputLabels = strings.Split(input_labels, ",")

	fmt.Println("start the server")
	ln, err := net.Listen("tcp", ingress_host+":"+ingress_port)
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
		go handleConnection(conn, egess_urls) // a goroutine handles conn so that the loop can accept other connections
	}
}
