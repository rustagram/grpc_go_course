package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/greet/greetpb"
    "io"
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

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
    result := ""
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&greetpb.LongGreetResponse{
                Result: result,
            })
        }
        if err != nil {
            return err
        }

        result += "Hello " + msg.GetGreeting().GetFirstName() + "! "
    }
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }

        if err != nil {
            return err
        }

        err = stream.Send(&greetpb.GreetEveryoneResponse{
            Result: "Hello " + req.GetGreeting().GetFirstName() + "! ",
        })
        if err != nil {
            return err
        }
    }
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
