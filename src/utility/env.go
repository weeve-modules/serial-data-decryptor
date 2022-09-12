package utility

import (
	"log"
	"os"
)

func GetEnvAsserted(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key, " was not found in the current environment")
	}
	return val
}
