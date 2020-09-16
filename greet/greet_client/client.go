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

    doServerStreaming(c)

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
