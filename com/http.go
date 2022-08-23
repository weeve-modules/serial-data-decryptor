package com

import (
	"bytes"
	"encoding/gob"
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
	log.Info("Starting the server")

	ln, err := net.Listen("tcp", utility.GetEnvAsserted("INGRESS_HOST")+":"+utility.GetEnvAsserted("INGRESS_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Started the server: %s", ln.Addr().String())

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
	err := dec.Decode(&msg)
	if err != nil {
		log.Error(err)
		return
	}

	if msg != nil {
		log.Error("Received message is nil")
		return
	}
	log.Infof("Received message: %+v\n", msg)

	data, err := processor.Process(msg)
	if err != nil {
		log.Error(err)
		return
	}
	conn.Close()

	//send processed data to egress endpoints
	sendData(data)
}

func sendData(data []byte) {
	egressUrls := strings.Split(strings.Replace(utility.GetEnvAsserted("EGRESS_URLS"), " ", "", -1), ",")
	for _, egressUrl := range egressUrls {
		log.Infof("Sending to: %s\n", egressUrl)

		r, err := http.NewRequest("POST", egressUrl, bytes.NewBuffer(data))
		if err != nil {
			log.Error(err)
			continue
		}

		r.Header.Add("Content-Type", "application/json")

		client := http.Client{}
		res, err := client.Do(r)
		if err != nil {
			log.Error(err)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Errorf("Sending failed with status code: %d\n", res.StatusCode)
			continue
		}

		rBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("Sent message response: %s\n", rBody)
	}
}
