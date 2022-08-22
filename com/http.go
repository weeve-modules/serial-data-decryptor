package com

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"serial-data-decryptor/models"
	"serial-data-decryptor/processor"
	"serial-data-decryptor/utility"
	"strings"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	fmt.Println("Starting the server")

	ln, err := net.Listen("tcp", utility.GetEnvAsserted("INGRESS_HOST")+":"+utility.GetEnvAsserted("INGRESS_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started the server: " + ln.Addr().String())

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			log.Error(err)
			continue
		}
		go handleMessages(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}

func handleMessages(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	msg := models.Data{}
	dec.Decode(&msg)
	fmt.Printf("Received message: %+v\n", msg)

	data, err := processor.Process(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()

	//send processed data to egress endpoints
	sendData(data)
}

func sendData(data []byte) {
	egress_Urls := strings.Split(utility.GetEnvAsserted("EGRESS_URLS"), ",")
	for _, egress_Url := range egress_Urls {
		fmt.Printf("Sending to: %s\n", egress_Url)

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

		if res.StatusCode != http.StatusOK {
			fmt.Printf("Sending failed with status code: %d\n", res.StatusCode)
			continue
		}

		rBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Sent message response: %s\n", rBody)
	}
}
