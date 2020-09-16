package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/calculator/calculatorpb"
    "google.golang.org/grpc"
    "io"
    "log"
)

func main() {
    fmt.Println("Hello I'm a client")

    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    c := calculatorpb.NewCalculatorServiceClient(conn)

    //doUnarySum(c)

    doServerStreaming(c)
}


func doUnarySum(c calculatorpb.CalculatorServiceClient) {
    resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
        Sum: &calculatorpb.Sum{
            A: 10,
            B: 3,
        },
    })
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    fmt.Println(resp.Sum)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
    resp, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PrimeNumberDecompositionRequest{
        Number: 120,
    })
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    for {
        msg, err := resp.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("could not connect: %v", err)
        }

        fmt.Println(msg)
    }
}
