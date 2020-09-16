package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/calculator/calculatorpb"
    "google.golang.org/grpc"
    "log"
    "net"
)

type server struct{}

func (s *server) PrimeNumberDecomposition(
    req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
    var (
        n, k int64
    )
    n = req.GetNumber()
    k = 2
    for n > 1 {
        if n%k == 0{
            err := stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
                PrimeNumber: k,
            })
            if err != nil {
                log.Fatalf("failed to send: %v", err)
            }
            n = n / k
        } else {
            k++
        }
    }

    return nil
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
    return &calculatorpb.SumResponse{
        Sum: req.GetSum().GetA() + req.GetSum().GetB(),
    }, nil
}

func main() {
    fmt.Println("Hello world")

    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("failed to listen: %v",err)
    }

    s := grpc.NewServer()

    calculatorpb.RegisterCalculatorServiceServer(s, &server{})

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
