package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetEnvAsserted(envVarName string) string {
	var thisEnvVar = os.Getenv(envVarName)
	if len(thisEnvVar) == 0 {
		log.Fatal(envVarName, " was not found in the current environment")
	}
	return thisEnvVar
}

type PlainFormatter struct {
	TimestampFormat string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	return []byte(fmt.Sprintf("%s %s : %s\n", timestamp, entry.Level, entry.Message)), nil
}

type P struct {
	M, N int64
}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &P{}
	dec.Decode(p)
	fmt.Printf("Received : %+v", p)
	conn.Close()
}

func main() {
	module_name := GetEnvAsserted("MODULE_NAME")
	ingress_host := GetEnvAsserted("INGRESS_HOST")
	ingress_port := GetEnvAsserted("INGRESS_PORT")
	egess_urls := GetEnvAsserted("EGRESS_URLS")

	log.Info("%s running on %s at port %s with end-point set to %s", module_name, ingress_host, ingress_port, egess_urls)

	fmt.Println("start the server")
	ln, err := net.Listen("tcp", ":"+ingress_port)
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
