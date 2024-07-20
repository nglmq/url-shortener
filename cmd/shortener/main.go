package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/server"
)

func main() {
	r, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(config.FlagRunAddr, r))
}
