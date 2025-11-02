package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	mediametav1 "github.com/AlexAnd012/mediameta/gen/go/mediameta/v1"
	"github.com/AlexAnd012/mediameta/internal/service"
	"github.com/AlexAnd012/mediameta/internal/storage"
)

func main() {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := storage.NewPostgresRepo(db)
	svc := service.NewMetadataService(repo)

	s := grpc.NewServer()
	mediametav1.RegisterMetadataServiceServer(s, svc)

	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":50051"
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC listening on", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
