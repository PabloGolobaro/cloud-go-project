package main

import (
	"cloud-go-project/cmd/gRPC"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:5051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	client := gRPC.NewKeyValueClient(conn)
	var action, key, value string
	if len(os.Args) > 2 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch action {
	case "get":
		responce, err := client.Get(ctx, &gRPC.GetRequest{Key: key})
		if err != nil {
			log.Fatal("could not get value for key %s: %v", key, err)
		}
		log.Printf("Get %s returns: %s", key, responce.Value)
	case "put":
		_, err := client.Put(ctx, &gRPC.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatal("could not put value for key %s: %v", key, err)
		}
		log.Printf("Put %s", key)
	case "delete":
		_, err := client.Delete(ctx, &gRPC.DeleteRequest{Key: key})
		if err != nil {
			log.Fatal("could not delte value for key %s: %v", key, err)
		}
		log.Printf("Delete %s", key)

	}

}
