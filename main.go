package main

import (
	"log"

	"github.com/YagoSchramm/Golinkr/infrastructure"
)

func main() {
	srv, cleanup, err := infrastructure.BuildServer()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	log.Fatal(srv.ListenAndServe())
}
