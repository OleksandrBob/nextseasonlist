package aws

import (
	"encoding/json"
	"log"
	"os"
)

func ResolveSecret(name string, region string) string {
	envSecret := os.Getenv(name)
	if envSecret == "" {
		return ""
	}

	if len(envSecret) > 0 && envSecret[0] != '{' {
		return envSecret
	}

	// Parse, in case value is taken from Secrets manager
	var secretParsed map[string]string
	err := json.Unmarshal([]byte(envSecret), &secretParsed)
	if err != nil {
		log.Fatal("Failed to parse secret JSON: ", err)
	}

	return secretParsed[name]
}
