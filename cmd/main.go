package main

import (
	"fmt"
	"serial-data-decryptor/com"
	"serial-data-decryptor/utility"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type PlainFormatter struct {
	TimestampFormat string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	return []byte(fmt.Sprintf("%s %s : %s\n", timestamp, entry.Level, entry.Message)), nil
}

func init() {
	const dateTimeFormat = "2006-01-02 15:04:05"

	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = dateTimeFormat
	log.SetFormatter(plainFormatter)
}

func main() {
	err := godotenv.Load("../docker/docker.env") // TODO: Only use it for testing locally
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s running on %s at port %s with end-point set to %s",
		utility.GetEnvAsserted("MODULE_NAME"),
		utility.GetEnvAsserted("INGRESS_HOST"),
		utility.GetEnvAsserted("INGRESS_PORT"),
		utility.GetEnvAsserted("EGRESS_URLS"))

	com.StartServer()
}
