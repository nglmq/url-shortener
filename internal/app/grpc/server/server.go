package grpcserver

import (
	"fmt"
	"log"
	"net"

	grpchandlers "github.com/nglmq/url-shortener/internal/app/grpc/handlers"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	pb "github.com/nglmq/url-shortener/internal/app/proto"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"github.com/nglmq/url-shortener/internal/app/storage/db"

	"github.com/nglmq/url-shortener/config"

	"google.golang.org/grpc"
)

// StartGRPCServer запускает gRPC сервер
func StartGRPCServer() error {
	store := storage.NewMemoryURLStore()
	shortener := &handlers.URLShortener{
		Store: store,
	}

	if config.DBConnection != "" {
		dbStorage, err := db.InitDBConnection()
		if err != nil {
			return err
		}
		shortener.DBStorage = dbStorage
	}

	if config.FlagInMemoryStorage != "" && config.DBConnection == "" {
		fileStore, err := storage.NewFileStorage(config.FlagInMemoryStorage)
		if err != nil {
			return err
		}
		shortener.FileStorage = fileStore

		if err = fileStore.ReadURLsFromFile(store.URLs); err != nil {
			log.Printf("Error reading URLs from file: %v", err)
			return err
		}
	}

	lis, err := net.Listen("tcp", ":3200")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterURLShortenerServer(grpcServer, &grpchandlers.URLShortenerServer{
		DBStorage: shortener.DBStorage,
		Store:     shortener.Store,
	})

	log.Printf("Starting gRPC server on :3200")
	return grpcServer.Serve(lis)
}
