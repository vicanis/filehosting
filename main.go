package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vicanis/filehosting/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error while reading .env: %s", err)
	}

	log.Fatal(server.Start())
}
