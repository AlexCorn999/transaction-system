package main

import (
	"log"

	"github.com/AlexCorn999/transaction-system/internal/transport"
)

func main() {

	server := transport.NewAPIServer()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
