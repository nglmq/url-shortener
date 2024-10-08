package main

import (
	"context"
	"errors"
	"fmt"
	grpcserver "github.com/nglmq/url-shortener/internal/app/grpc/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/nglmq/url-shortener/internal/app/cert"

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

	if config.GRPCServer {
		if err := grpcserver.StartGRPCServer(); err != nil {
			log.Fatal(err)
		}
	} else {
		srv := &http.Server{
			Addr:    config.FlagRunAddr,
			Handler: r,
		}

		idleConnsClosed := make(chan struct{})
		sigint := make(chan os.Signal, 3)

		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		go func() {
			<-sigint

			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
			close(idleConnsClosed)
		}()

		if config.EnableHTTPS {
			cert.CertGen()
			fmt.Printf("Starting server: %s", config.FlagRunAddr)
			err = srv.ListenAndServeTLS("cert.pem", "key.pem")
		} else {
			err = srv.ListenAndServe()
		}

		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}

		<-idleConnsClosed
		fmt.Println("Server Shutdown gracefully")
	}
}
