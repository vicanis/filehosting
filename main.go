package main

import (
	"log"

	"github.com/vicanis/filehosting/server"
)

func main() {
	log.Fatal(server.Start())
}
