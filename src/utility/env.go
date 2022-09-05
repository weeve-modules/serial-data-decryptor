package utility

import (
	"log"
	"os"
)

func GetEnvAsserted(envVarName string) string {
	val, ok := os.LookupEnv(envVarName)
	if !ok {
		log.Fatal(envVarName, " was not found in the current environment")
	}
	return val
}
