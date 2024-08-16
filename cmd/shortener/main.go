package main

import (
	"fmt"
	"github.com/nglmq/url-shortener/internal/app/cert"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	r, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	if config.EnableHTTPS {
		cert.CertGen()

		log.Fatal(http.ListenAndServeTLS(config.FlagRunAddr,
			"cert.pem",
			"key.pem",
			r))
	}

	log.Fatal(http.ListenAndServe(config.FlagRunAddr, r))
}
