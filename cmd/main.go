package main

import (
	"log"

	"github.com/Asif-Faizal/Gommerce/cmd/api"
)

func main() {
	server := api.NewAPIServer(":3000", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
