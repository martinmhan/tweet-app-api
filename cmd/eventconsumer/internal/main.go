package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/application"
	"github.com/martinmhan/tweet-app-api/cmd/eventconsumer/internal/domain/rpcclient"
)

func main() {
	godotenv.Load()

	mqhost := os.Getenv("MQ_HOST")
	mqport := os.Getenv("MQ_PORT")
	mqname := os.Getenv("MQ_NAME")
	dahost := os.Getenv("DA_HOST")
	daport := os.Getenv("DA_PORT")
	rvhost := os.Getenv("RV_HOST")
	rvport := os.Getenv("RV_PORT")
	if mqport == "" || mqhost == "" || mqname == "" || dahost == "" || daport == "" || rvhost == "" || rvport == "" {
		log.Fatal("Missing environment variable(s). Please edit .env file")
	}

	da := rpcclient.DatabaseAccess{
		Host: dahost,
		Port: daport,
	}
	err := da.Connect()
	if err != nil {
		log.Fatal("Failed to connect to Database Access server")
	}

	rv := rpcclient.ReadView{
		Host: rvhost,
		Port: rvport,
	}
	err = rv.Connect()
	if err != nil {
		log.Fatal("Failed to connect to Read View server")
	}

	defer da.Disconnect()
	defer rv.Disconnect()

	ec := &application.EventConsumer{
		DatabaseAccess:   da,
		ReadView:         rv,
		MessageQueueHost: mqhost,
		MessageQueuePort: mqport,
		MessageQueueName: mqname,
	}

	ec.Init()
}
