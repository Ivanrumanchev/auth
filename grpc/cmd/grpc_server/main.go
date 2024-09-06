package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Ivanrumanchev/auth/grpc/pkg/user_v1"
)

func generateChatID() int64 {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return t
}

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create Name: %s", req.GetName())
	log.Printf("Create Email: %s", req.GetEmail())
	log.Printf("Create Password: %s", req.GetPassword())
	log.Printf("Create PasswordConfirm: %s", req.GetPasswordConfirm())
	log.Printf("Create Role: %s", req.GetRole())

	return &desc.CreateResponse{
		Id: generateChatID(),
	}, nil
}

// Update
func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Update Name: %s", req.GetName())
	log.Printf("Update Email: %s", req.GetEmail())
	log.Printf("Update Role: %s", req.GetRole())
	log.Printf("Update Id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

// Get
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get UserId: %d", req.GetId())

	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.BeerName(),
		Email:     gofakeit.IPv4Address(),
		Role:      desc.Role_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Delete
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Delete User id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
