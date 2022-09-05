package com

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"serial-data-decryptor/models"
	"serial-data-decryptor/processor"
	"serial-data-decryptor/utility"
	"strings"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	log.Info("Starting the server")

	addr := utility.GetEnvAsserted("INGRESS_HOST") + ":" + utility.GetEnvAsserted("INGRESS_PORT")

	http.HandleFunc("/", handleMessages)

	log.Fatal(http.ListenAndServe(addr, nil))

	log.Info("Started the server on ", addr)
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Received message: ", string(body))

	var data models.Data

	json.Unmarshal(body, &data)

	resp, err := processor.Process(data)
	if err != nil {
		log.Error(err)
		return
	}

	//send processed data to egress endpoints
	sendData(resp)
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

		rBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("Sent message response: %s\n", rBody)
	}
}
