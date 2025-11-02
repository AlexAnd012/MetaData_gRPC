package service

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mediametav1 "github.com/AlexAnd012/mediameta/gen/go/mediameta/v1"
	"github.com/AlexAnd012/mediameta/internal/storage"
)

type MetadataService struct {
	mediametav1.UnimplementedMetadataServiceServer
	repo storage.Repository
}

func NewMetadataService(r storage.Repository) *MetadataService { return &MetadataService{repo: r} }

func (s *MetadataService) HealthCheck(context.Context, *mediametav1.HealthCheckRequest) (*mediametav1.HealthCheckResponse, error) {
	return &mediametav1.HealthCheckResponse{Status: "ok"}, nil
}

func (s *MetadataService) CreateMetadata(ctx context.Context, req *mediametav1.CreateMetadataRequest) (*mediametav1.CreateMetadataResponse, error) {
	if req.GetFilename() == "" {
		return nil, status.Error(codes.InvalidArgument, "filename required")
	}
	now := time.Now().Unix()
	m := &mediametav1.Metadata{
		Id:          uuid.NewString(),
		Filename:    req.GetFilename(),
		SizeBytes:   req.GetSizeBytes(),
		ContentType: req.GetContentType(),
		OwnerId:     req.GetOwnerId(),
		CreatedAt:   now, //в Unix-секундах
	}
	if err := s.repo.Insert(ctx, m); err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &mediametav1.CreateMetadataResponse{Meta: m}, nil
}

func (s *MetadataService) CreateMetadataFromPath(ctx context.Context, req *mediametav1.CreateFromPathRequest) (*mediametav1.CreateFromPathResponse, error) {
	p := req.GetPath()
	if p == "" {
		return nil, status.Error(codes.InvalidArgument, "path required")
	}
	fi, err := os.Stat(p) // берём размер
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "stat: %v", err)
	}
	filename := filepath.Base(p)
	ct := mime.TypeByExtension(filepath.Ext(filename)) // смотрим расширение файла
	if ct == "" {
		ct = "application/octet-stream" // просто двоичные данные
	}
	now := time.Now().Unix()
	m := &mediametav1.Metadata{
		Id:          uuid.NewString(),
		Filename:    filename,
		SizeBytes:   fi.Size(),
		ContentType: ct,
		OwnerId:     req.GetOwnerId(),
		CreatedAt:   now,
	}
	if err := s.repo.Insert(ctx, m); err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &mediametav1.CreateFromPathResponse{Meta: m}, nil
}

func (s *MetadataService) GetMetadata(ctx context.Context, req *mediametav1.GetMetadataRequest) (*mediametav1.GetMetadataResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}
	m, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return &mediametav1.GetMetadataResponse{Meta: m}, nil
}

func (s *MetadataService) ListMetadata(ctx context.Context, req *mediametav1.ListMetadataRequest) (*mediametav1.ListMetadataResponse, error) {
	ps := req.GetPageSize()
	if ps <= 0 || ps > 100 {
		ps = 20
	}
	offset := 0
	if tok := req.GetPageToken(); tok != "" {
		_, _ = fmt.Sscanf(tok, "%d", &offset)
	}
	items, next, err := s.repo.List(ctx, int(ps), offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list: %v", err)
	}
	token := ""
	if next > 0 {
		token = fmt.Sprintf("%d", next)
	}
	return &mediametav1.ListMetadataResponse{Items: items, NextPageToken: token}, nil
}

func (s *MetadataService) UpdateMetadata(ctx context.Context, req *mediametav1.UpdateMetadataRequest) (*mediametav1.UpdateMetadataResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}
	cur, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	if req.GetFilename() != "" {
		cur.Filename = req.GetFilename()
	}
	if req.GetContentType() != "" {
		cur.ContentType = req.GetContentType()
	}
	if err := s.repo.Update(ctx, cur); err != nil {
		return nil, status.Errorf(codes.Internal, "update: %v", err)
	}
	return &mediametav1.UpdateMetadataResponse{Meta: cur}, nil
}

func (s *MetadataService) DeleteMetadata(ctx context.Context, req *mediametav1.DeleteMetadataRequest) (*mediametav1.DeleteMetadataResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	return &mediametav1.DeleteMetadataResponse{}, nil
}
