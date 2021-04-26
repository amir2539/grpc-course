package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Blog Client Started")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connet %v", err)
	}

	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)

	//createBlog(c)
	readBlog(c)
}
func readBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Reading blog")

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "6086ffdd4e6d91ce1488d041"})
	if err != nil {
		log.Printf("Error %v", err)
		return
	}

	fmt.Printf("Blog is:%v", res.GetBlog())

}

func createBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "One",
		Content:  "First content",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected eror %v", err)
	}

	fmt.Println("Blog has been created %v", res)
}
