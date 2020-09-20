package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/calculator/calculatorpb"
    "google.golang.org/grpc"
    "io"
    "log"
    "time"
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

    //doServerStreaming(c)

    //doClientStreaming(c)

    doBiDiStreamig(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
    reqs := []*calculatorpb.ComputeAverageRequest{
        {
            Number: 1,
        },
        {
            Number: 2,
        },
        {
            Number: 3,
        },
        {
            Number: 4,
        },
    }

    stream, err := c.ComputeAverage(context.Background())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    for _, req := range reqs {
        err := stream.Send(req)
        if err != nil {
            log.Fatalf("could not connect: %v", err)
        }
    }

    resp, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    fmt.Println(resp)
}

func doBiDiStreamig(c calculatorpb.CalculatorServiceClient) {
    reqs := []*calculatorpb.FindMaximumRequest{
        {
            Number: 1,
        },
        {
            Number: 5,
        },
        {
            Number: 3,
        },
        {
            Number: 6,
        },
        {
            Number: 2,
        },
        {
            Number: 20,
        },
    }

    stream, err := c.FindMaximum(context.Background())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    waitch := make(chan struct{})

    go func() {
        for _, req := range reqs {
            err = stream.Send(&calculatorpb.FindMaximumRequest{
                Number: req.GetNumber(),
            })
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatalf("could not connect: %v", err)
            }
            time.Sleep(time.Second)
        }
        err := stream.CloseSend()
        if err != nil {
            log.Fatalf("could not connect: %v", err)
        }
    }()

    go func() {
        for {
            resp, err := stream.Recv()
            if err != nil {
                log.Fatalf("could not connect: %v", err)
            }

            fmt.Println(resp.CurrentMaximum)
        }
    }()

    <-waitch
}
