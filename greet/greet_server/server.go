package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/greet/greetpb"
    "log"
    "net"
    "strconv"

    "google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
    return &greetpb.GreetResponse{
        Result: "hello " + req.GetGreeting().GetFirstName(),
    }, nil
}


func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error{
    for i := 0; i < 10; i++ {
        err := stream.Send(&greetpb.GreetManyTimesResponse{
            Result: "hello " + req.GetGreeting().GetFirstName() + " number " + strconv.Itoa(i),
        })
        if err != nil {
            return err
        }
    }
    return nil
}


func main() {
    fmt.Println("Hello world")

    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("failed to listen: %v",err)
    }

    s := grpc.NewServer()

    greetpb.RegisterGreetServiceServer(s, &server{})

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
