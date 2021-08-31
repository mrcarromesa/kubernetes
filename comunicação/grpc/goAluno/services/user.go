package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mrcarromesa/grpc/pb"
)

// type UserServiceClient interface {
// 	AddUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
// 	AddUserVerbose(ctx context.Context, in *User, opts ...grpc.CallOption) (UserService_AddUserVerboseClient, error)
// 	AddUsers(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUsersClient, error)
// 	AddUserStreamBoth(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUserStreamBothClient, error)
// }

type UserService struct {
	pb.UnimplementedUserServiceServer // colocar essa parada aqui e boas pois se eu não implementar algo que esteja no proto não dará error!
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	// Insert - Database

	fmt.Println(req.Name)

	return &pb.User{
		Id:    "123",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetEmail(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetEmail(),
			Email: req.GetEmail(),
		},
	})

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {

	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding", req.GetName())

	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending stream to the client: %v", err)
		}
	}
}
