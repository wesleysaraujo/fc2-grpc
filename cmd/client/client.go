package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/wesleysaraujo/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect to gRPC server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUsersStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jonas",
		Email: "jonas@gmail.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatal("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jonas",
		Email: "jonas@gmail.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatal("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Coud not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	users := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Wesley Serafim",
			Email: "wesley@agits.com.br",
		},
		&pb.User{
			Id:    "2",
			Name:  "Manuella Alves Serafim de Araújo",
			Email: "manu@agits.com.br",
		},
		&pb.User{
			Id:    "3",
			Name:  "Jessica Alves",
			Email: "jessica@agits.com.br",
		},
		&pb.User{
			Id:    "4",
			Name:  "Kaua Alves",
			Email: "kaua@agits.com.br",
		},
		&pb.User{
			Id:    "5",
			Name:  "Henrique Alves",
			Email: "henrique@agits.com.br",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range users {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUsersStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	users := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Wesley Serafim",
			Email: "wesley@agits.com.br",
		},
		&pb.User{
			Id:    "2",
			Name:  "Manuella",
			Email: "manu@agits.com.br",
		},
		&pb.User{
			Id:    "3",
			Name:  "Jessica",
			Email: "jessica@agits.com.br",
		},
		&pb.User{
			Id:    "4",
			Name:  "Kaua BR",
			Email: "kaua@agits.com.br",
		},
		&pb.User{
			Id:    "5",
			Name:  "Henrique A",
			Email: "henrique@agits.com.br",
		},
	}

	wait := make(chan int)

	// Goroutine - Sending data to the server
	go func() {
		for _, req := range users {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	// Goroutine - Receiving data from the server
	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}

			fmt.Printf("Receiving user %v with status %v \n", res.GetUser().GetName(), res.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
