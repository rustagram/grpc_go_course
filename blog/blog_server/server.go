package main

import (
    "context"
    "fmt"
    "github.com/rustagram/grpc-go-course/blog/blogpb"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "log"
    "net"
    "os"
    "os/signal"
    "time"
)

var (
    collection *mongo.Collection
)

type server struct {
}

type blogItem struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    AuthorID string `bson:"author_id,omitempty"`
    Content string `bson:"content,omitempty"`
    Title string `bson:"title,omitempty"`
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
    data := blogItem{
        AuthorID: req.GetBlog().GetAuthorId(),
        Title: req.GetBlog().GetTitle(),
        Content: req.GetBlog().GetContent(),
    }

    res, err := collection.InsertOne(context.Background(), data)
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    oid, ok := res.InsertedID.(primitive.ObjectID)
    if !ok {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID"))
    }

    return &blogpb.CreateBlogResponse{
        Blog: &blogpb.Blog{
            Id: oid.Hex(),
            AuthorId: req.GetBlog().GetAuthorId(),
            Title: req.GetBlog().GetTitle(),
            Content: req.GetBlog().GetContent(),
        },
    }, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
    oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    blog := blogItem{}
    filter := bson.D{{"_id", oid}}

    err = collection.FindOne(context.Background(), filter).Decode(&blog)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Document not found: %v", err))
        }
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    return &blogpb.ReadBlogResponse{
        Blog: dataToBlogPb(blog),
    }, nil
}

func dataToBlogPb(data blogItem) *blogpb.Blog {
    return &blogpb.Blog{
        Id: data.ID.Hex(),
        AuthorId: data.AuthorID,
        Title: data.Title,
        Content: data.Content,
    }
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
    oid, err := primitive.ObjectIDFromHex(req.GetBlog().GetId())
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    data := blogItem{
        AuthorID: req.GetBlog().GetAuthorId(),
        Content: req.GetBlog().GetContent(),
        Title: req.GetBlog().GetTitle(),
    }
    filter := bson.D{{"_id", oid}}

    update := bson.M{
        "$set": data,
    }


    _, err = collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    return &blogpb.UpdateBlogResponse{
        Blog: dataToBlogPb(data),
    }, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
    oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    filter := bson.D{{"_id", oid}}
    res, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    if res.DeletedCount == 0 {
        return nil, status.Errorf(codes.Internal, fmt.Sprintf("Couldn't found blog with this id"))
    }

    return &blogpb.DeleteBlogResponse{
        BlogId: req.GetBlogId(),
    }, nil
}

func (s *server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
    fmt.Println("List blog request")

    cursor, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        return status.Errorf(codes.Internal, fmt.Sprintf("failed to list blogs: %v", err))
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        data := blogItem{}
        err := cursor.Decode(&data)
        if err != nil {
            return status.Errorf(codes.Internal, fmt.Sprintf("Couldn't decode: %v", err))
        }

        err = stream.Send(&blogpb.ListBlogResponse{
            Blog: dataToBlogPb(data),
        })
        if err != nil {
            return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
        }
    }
    if err = cursor.Err(); err != nil {
        return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
    }

    return nil
}

func main() {
    fmt.Println("starting blog service")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect to mongodb
    fmt.Println("Connecting to MongoDB")
    client, err := mongo.Connect(ctx,
        options.Client().ApplyURI("mongodb://localhost:27017"))

    defer func() {
        if err = client.Disconnect(ctx); err != nil {
            log.Fatalf("failed to disconnect: %v",err)
        }
    }()

    collection = client.Database("mydb").Collection("blog")

    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("failed to listen: %v",err)
    }

    s := grpc.NewServer()

    blogpb.RegisterBlogServiceServer(s, &server{})

    go func() {
        fmt.Println("starting server")
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // Wait for control C to exit
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt)

    // Block until a signal is received
    <-ch
    fmt.Println("Stopping the server")
    s.Stop()
    fmt.Println("Closing the listener")
    lis.Close()
    fmt.Println("Closing MongoDB Connection")
    client.Disconnect(ctx)
    fmt.Println("End of program")


}

