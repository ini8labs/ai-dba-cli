package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/ini8labs/ai-dba-cli/cmd"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		logrus.Fatal("WEBHOOK_URL environment variable is not set")
	}

	cmd.Execute()
}
