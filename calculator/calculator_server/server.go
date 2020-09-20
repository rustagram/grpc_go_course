package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/calculator/calculatorpb"
    "google.golang.org/grpc"
    "io"
    "log"
    "net"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
    return &calculatorpb.SumResponse{
        Sum: req.GetSum().GetA() + req.GetSum().GetB(),
    }, nil
}

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

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
    var sum, count int64
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
                Average: float64(sum)/float64(count),
            })
        }
        if err != nil {
            return err
        }

        sum += msg.Number
        count++
    }
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
    var max int64 = 0
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        if req.GetNumber() > max {
            max = req.GetNumber()
            if err := stream.Send(&calculatorpb.FindMaximumResponse{
                CurrentMaximum: max,
            }); err != nil {
                return err
            }
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

    calculatorpb.RegisterCalculatorServiceServer(s, &server{})

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
