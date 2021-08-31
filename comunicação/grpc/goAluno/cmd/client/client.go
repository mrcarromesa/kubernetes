package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mrcarromesa/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close() // o defer ver quando a variavel parou de ser utilizado e dai fecha a conexao

	client := pb.NewUserServiceClient(connection)

	// AddUser(client) // <- Comentado para utilizar o stream

	// AddUserVerbose(client)

	// AddUsers(client)

	AddUserStreamBoth(client) // <- bi-direcional
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)

}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status:", stream.Status, "-", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "R1",
			Name:  "Rodolfo",
			Email: "example1@email.com",
		},
		{
			Id:    "R2",
			Name:  "Rodolfo 2",
			Email: "example2@email.com",
		},
		{
			Id:    "R3",
			Name:  "Rodolfo 3",
			Email: "example3@email.com",
		},
		{
			Id:    "R4",
			Name:  "Rodolfo 4",
			Email: "example4@email.com",
		},
		{
			Id:    "R5",
			Name:  "Rodolfo 5",
			Email: "example5@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "R1",
			Name:  "Rodolfo",
			Email: "example1@email.com",
		},
		{
			Id:    "R2",
			Name:  "Rodolfo 2",
			Email: "example2@email.com",
		},
		{
			Id:    "R3",
			Name:  "Rodolfo 3",
			Email: "example3@email.com",
		},
		{
			Id:    "R4",
			Name:  "Rodolfo 4",
			Email: "example4@email.com",
		},
		{
			Id:    "R5",
			Name:  "Rodolfo 5",
			Email: "example5@email.com",
		},
	}

	wait := make(chan int)

	// Para ficar enviando e
	// anonymos func
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		// Parei de enviar aqui
		stream.CloseSend()
	}()

	// Em paralelo/concorrente
	// ficar recebendo
	// Quando o servidor para de enviar ele cai no break e é encerrada
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Receiving user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
