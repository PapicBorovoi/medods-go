package main

import (
	"fmt"
	"net/http"

	"github.com/PapicBorovoi/medods-go/internals/db"
	"github.com/PapicBorovoi/medods-go/internals/handler"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const addr string = "localhost:3000"

func main() {
	log.SetReportCaller(true)
	godotenvErr := godotenv.Load()

	if godotenvErr != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Server is starting...")

	mux := http.NewServeMux()
	handler.Handle(mux)

	dbClient, dbErr := db.Connect()

	if dbErr != nil {
		log.Fatal(dbErr)
	}


	fmt.Printf(`Server ready and listening on %s`, addr)

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		db.Close(dbClient)
		log.Fatal(err)
	}
}