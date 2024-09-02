package grpchandlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nglmq/url-shortener/internal/app/proto"
	"github.com/nglmq/url-shortener/internal/app/storage/db"

	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"github.com/nglmq/url-shortener/internal/app/storage"
)

// URLShortenerServer
type URLShortenerServer struct {
	pb.UnimplementedURLShortenerServer
	DBStorage *db.PostgresStorage
	Store     storage.URLStore
}

// SaveURL
func (s *URLShortenerServer) SaveURL(ctx context.Context, req *pb.SaveURLRequest) (*pb.SaveURLResponse, error) {
	alias := random.NewRandomURL()
	userID := req.GetUserId()

	existAlias, err := s.DBStorage.SaveURL(ctx, userID, alias, req.GetUrl())
	if err != nil {
		return nil, fmt.Errorf("failed to save URL to database: %v", err)
	}

	shortURL := fmt.Sprintf(config.FlagBaseURL + "/" + existAlias)
	return &pb.SaveURLResponse{ShortUrl: shortURL}, nil
}

// SaveJSON
func (s *URLShortenerServer) SaveJSON(ctx context.Context, req *pb.SaveJSONRequest) (*pb.SaveJSONResponse, error) {
	alias := random.NewRandomURL()
	userID := req.GetUserId()

	existAlias, err := s.DBStorage.SaveURL(ctx, userID, alias, req.GetJsonUrl())
	if err != nil {
		return nil, fmt.Errorf("failed to save URL to database: %v", err)
	}

	shortURL := config.FlagBaseURL + "/" + existAlias
	return &pb.SaveJSONResponse{ShortUrl: shortURL}, nil
}

// SaveJSONBatch
func (s *URLShortenerServer) SaveJSONBatch(ctx context.Context, req *pb.SaveJSONBatchRequest) (*pb.SaveJSONBatchResponse, error) {
	var results []*pb.BatchURLResponse

	for _, batchReq := range req.GetUrls() {
		alias := random.NewRandomURL()

		existAlias, err := s.DBStorage.SaveURL(ctx, req.GetUserId(), alias, batchReq.GetOriginalUrl())
		if err != nil {
			return nil, fmt.Errorf("failed to save URL to database: %v", err)
		}

		shortURL := config.FlagBaseURL + "/" + existAlias
		results = append(results, &pb.BatchURLResponse{
			CorrelationId: batchReq.GetCorrelationId(),
			ShortUrl:      shortURL,
		})
	}

	return &pb.SaveJSONBatchResponse{Results: results}, nil
}

// GetURL
func (s *URLShortenerServer) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	url, deleted, err := s.DBStorage.GetURL(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get URL: %v", err)
	}

	return &pb.GetURLResponse{Url: url, Deleted: deleted}, nil
}

// GetAllURLs
func (s *URLShortenerServer) GetAllURLs(ctx context.Context, req *pb.GetAllURLsRequest) (*pb.GetAllURLsResponse, error) {
	urls, err := s.DBStorage.GetAllUserURLs(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %v", err)
	}

	var urlMappings []*pb.URLMapping
	for alias, originalURL := range urls {
		urlMappings = append(urlMappings, &pb.URLMapping{
			ShortUrl:    config.FlagBaseURL + "/" + alias,
			OriginalUrl: originalURL,
		})
	}

	return &pb.GetAllURLsResponse{Urls: urlMappings}, nil
}

// GetStats
func (s *URLShortenerServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	uriQuantity, usersQuantity, err := s.DBStorage.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}

	return &pb.GetStatsResponse{Urls: int32(uriQuantity), Users: int32(usersQuantity)}, nil
}

// DeleteURL
func (s *URLShortenerServer) DeleteURL(ctx context.Context, req *pb.DeleteURLRequest) (*pb.DeleteURLResponse, error) {
	err := s.DBStorage.DeleteURL(ctx, req.GetUserId(), req.GetAlias())
	if err != nil {
		return nil, fmt.Errorf("failed to delete URL: %v", err)
	}

	return &pb.DeleteURLResponse{Deleted: true}, nil
}
