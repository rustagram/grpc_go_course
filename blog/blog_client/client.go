package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/blog/blogpb"
    "google.golang.org/grpc"
    "io"
    "log"
)

func main() {
    fmt.Println("Blog client")

    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    c := blogpb.NewBlogServiceClient(conn)

    // Create Blog
    //fmt.Println("Creating Blog")
    //blog := &blogpb.Blog{
    //   AuthorId: "Rustam",
    //   Title: "My first blog",
    //   Content: "Content of my first blog",
    //}

    //createdBlog, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
    //   Blog: blog,
    //})
    //if err != nil {
    //   log.Fatalf("Unexpected error %v", err.Error())
    //}
    //fmt.Printf("Blog has been created %v", createdBlog)

    //readBlogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
    //    BlogId: "5fc3ae164235f0861ca13231",
    //})
    //if err != nil {
    //    log.Fatalf("Unexpected error %v", err.Error())
    //}
    //
    //fmt.Println(readBlogRes.Blog)

    //updatedBlog, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
    //    Blog: &blogpb.Blog{
    //        Id: "5fc3ae164235f0861ca13230",
    //        AuthorId: "Rustam (edited)",
    //        Title: "My second blog (edited)",
    //        Content: "Content of my first blog (edited)",
    //    },
    //})
    //if err != nil {
    //   log.Fatalf("Unexpected error %v", err.Error())
    //}
    //
    //fmt.Println(updatedBlog)

    //_, err = c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
    //    BlogId: "5fc3ae164235f0861ca13230",
    //})
    //if err != nil {
    //  log.Fatalf("Unexpected error %v", err.Error())
    //}

    resp, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
    if err != nil {
     log.Fatalf("Unexpected error %v", err.Error())
    }

    for {
        msg, err := resp.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("could not connect: %v", err)
        }

        fmt.Println(msg.Blog)
    }
}


