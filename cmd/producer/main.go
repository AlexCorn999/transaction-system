package main

import (
	"log"

	"github.com/AlexCorn999/transaction-system/internal/producer"
)

func main() {
	server := producer.NewAPIServer()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
