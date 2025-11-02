package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mediametav1 "github.com/AlexAnd012/mediameta/gen/go/mediameta/v1"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: client <path> [owner_id]")
	}
	path := os.Args[1]
	owner := "u1"
	if len(os.Args) > 2 {
		owner = os.Args[2]
	}

	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := mediametav1.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.CreateMetadataFromPath(ctx, &mediametav1.CreateFromPathRequest{
		Path:    path,
		OwnerId: owner,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created id=%s name=%s size=%d ct=%s owner=%s\n",
		resp.Meta.Id, resp.Meta.Filename, resp.Meta.SizeBytes, resp.Meta.ContentType, resp.Meta.OwnerId)
}
