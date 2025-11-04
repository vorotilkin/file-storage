package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	putTTL = time.Minute * 10
	getTTL = time.Hour * 1
)

type Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type Service struct {
	client *minio.Client
	config *Config
}

func (s *Service) PresignPut(ctx context.Context, objectName string) (string, error) {
	url, err := s.client.PresignedPutObject(ctx, s.config.Bucket, objectName, putTTL)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *Service) PresignGet(ctx context.Context, objectName string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.config.Bucket, objectName, getTTL, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *Service) Bucket() string {
	return s.config.Bucket
}

func New(cfg Config) (*Service, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
		config: &cfg,
	}, nil
}
