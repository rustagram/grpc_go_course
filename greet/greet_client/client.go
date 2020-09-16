package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/greet/greetpb"
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

    c := greetpb.NewGreetServiceClient(conn)

    //doUnary(c)

    //doServerStreaming(c)

    doClientStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
    resp, err := c.Greet(context.Background(), &greetpb.GreetRequest{
        Greeting: &greetpb.Greeting{
            FirstName: "Rustam",
            LastName: "Turgunov",
        },
    })
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    fmt.Println(resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
    resp, err := c.GreetManyTimes(context.Background(), &greetpb.GreetManyTimesRequest{
        Greeting: &greetpb.Greeting{
            FirstName: "Rustam",
            LastName: "Turgunov",
        },
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

func doClientStreaming(c greetpb.GreetServiceClient) {
    stream, err := c.LongGreet(context.Background())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }

    reqs := []*greetpb.LongGreetRequest{
        {
            Greeting: &greetpb.Greeting{
                FirstName: "Rustam",
            },
        },
        {
            Greeting: &greetpb.Greeting{
                FirstName: "Akbar",
            },
        },
        {
            Greeting: &greetpb.Greeting{
                FirstName: "Begzod",
            },
        },
        {
            Greeting: &greetpb.Greeting{
                FirstName: "Tulkin",
            },
        },
        {
            Greeting: &greetpb.Greeting{
                FirstName: "Farhod",
            },
        },
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
