package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/Ivanrumanchev/auth/pkg/user_v1"
)

const (
	address = "localhost:50051"
	userID  = int64(4)
)

func closeConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer closeConn(conn)

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pass := gofakeit.FirstName()

	createResponse, err := c.Create(ctx, &desc.CreateRequest{
		Name:            gofakeit.BeerName(),
		Email:           gofakeit.IPv4Address(),
		Role:            desc.Role_ADMIN,
		Password:        pass,
		PasswordConfirm: pass,
	})
	if err != nil {
		log.Fatalf("failed to create user by id: %v", err)
	}

	log.Printf(color.RedString("Create User info:\n"), color.GreenString("%+v", createResponse.GetId()))

	deleteResponse, err := c.Delete(ctx, &desc.DeleteRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to delete user by id: %v", err)
	}

	log.Printf(color.RedString("Delete User info:\n"), color.GreenString("%+v", deleteResponse))

	getResponse, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", getResponse))

	updateResponse, err := c.Update(ctx, &desc.UpdateRequest{Id: int64(4), Role: desc.Role_USER})
	if err != nil {
		log.Fatalf("failed to update user by id: %v", err)
	}

	log.Printf(color.RedString("User update:\n"), color.GreenString("%+v", updateResponse))
}
