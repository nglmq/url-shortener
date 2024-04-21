package main

import (
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/server"
	"log"
	"net/http"
)

func main() {
	r, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(config.FlagRunAddr, r))
}
