package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func dataToBlogpb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorId: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal Error %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot parse to oid %v", err))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("invalid argument"))
	}

	data := &blogItem{}
	filter := bson.D{
		bson.E{Key: "_id", Value: oid},
	}
	result := collection.FindOne(context.Background(), filter)
	if err := result.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("CAnnot find blog"))
	}

	return &blogpb.ReadBlogResponse{Blog: dataToBlogpb(data)}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	fmt.Println("Update blog started")
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("invalid argument"))
	}

	data := &blogItem{}
	filter := bson.D{
		bson.E{Key: "_id", Value: oid},
	}

	result := collection.FindOne(context.Background(), filter)
	if err := result.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("CAnnot find blog"))
	}

	data.AuthorId = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, UpdateError := collection.ReplaceOne(context.Background(), filter, data)
	if UpdateError != nil {
		return nil, status.Errorf(
			codes.Internal, fmt.Sprintf("Cannot update %v", UpdateError))
	}

	return &blogpb.UpdateBlogResponse{Blog: dataToBlogpb(data)}, nil

}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("invalid argument"))
	}

	filter := bson.D{
		bson.E{Key: "_id", Value: oid},
	}
	deleteResult, deleteError := collection.DeleteOne(context.Background(), filter)
	if deleteError != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cant delete blog %v", deleteError))
	}

	if deleteResult.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog"))
	}

	return &blogpb.DeleteBlogResponse{}, nil

}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	fmt.Println("List Blog")

	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Cant adad blog %v", err))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &blogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Cant decode blog %v", err))
		}

		stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogpb(data)})
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("unknown %v", err))
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Server Started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	collection = client.Database("grpc").Collection("grpc")

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)
	go func() {
		fmt.Println("start server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to listen %v", err)
		}
	}()

	// Wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("Stopped")
	s.Stop()
	lis.Close()
}
