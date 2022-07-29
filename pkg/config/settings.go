package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

var EmailName = goDotEnvVariable("EMAIL_NAME")
var EmailPass = goDotEnvVariable("EMAIL_PASSWORD")
var EmailHost = goDotEnvVariable("EMAIL_HOST")
var EmailPort, _ = strconv.Atoi(goDotEnvVariable("EMAIL_PORT"))

var EmailsStoragePath = "storage/emails.json"
var BitcoinCoingateDomain = "https://api.coingate.com"
