package frontend

import (
	"cloud-go-project/hexarch/core"
	"cloud-go-project/hexarch/frontend/gRPC"
	"context"
	grpc "google.golang.org/grpc"
	"log"
	"net"
)

type grpcFrontEnd struct {
	store *core.KeyValueStore
	gRPC.UnimplementedKeyValueServer
}

func (g *grpcFrontEnd) Start(store *core.KeyValueStore) error {
	g.store = store
	s := grpc.NewServer()
	gRPC.RegisterKeyValueServer(s, g)
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	return s.Serve(lis)
}

func (s *grpcFrontEnd) Get(ctx context.Context, r *gRPC.GetRequest) (*gRPC.GetResponce, error) {
	log.Printf("Received GET key=%v", r.Key)
	value, err := s.store.Get(r.Key)
	return &gRPC.GetResponce{Value: value}, err
}
func (s *grpcFrontEnd) Put(ctx context.Context, r *gRPC.PutRequest) (*gRPC.PutResponce, error) {
	log.Printf("Received PUT key=%v, value=%v", r.Key, r.Value)
	err := s.store.Put(r.Key, r.Value)
	return &gRPC.PutResponce{}, err
}
func (s *grpcFrontEnd) Delete(ctx context.Context, r *gRPC.DeleteRequest) (*gRPC.DeleteResponce, error) {
	log.Printf("Received DELETE key=%v", r.Key)
	err := s.store.Delete(r.Key)
	return &gRPC.DeleteResponce{}, err

}
