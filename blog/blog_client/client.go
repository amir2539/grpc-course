package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/blog/blogpb"
	"google.golang.org/grpc"
	"io"
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
	//readBlog(c)
	//updateBlog(c)
	//deleteBlog(c)
	listBlogs(c)
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
		AuthorId: "hgfb",
		Title:    "tr3",
		Content:  "سلام ایرانی",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected eror %v", err)
	}

	fmt.Println("Blog has been created %v", res)
}

func updateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		Id:       "6086ffdd4e6d91ce1488d041",
		AuthorId: "author changed",
		Title:    "One (edited)",
		Content:  "First content 2",
	}
	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected eror %v", err)
	}

	fmt.Println("Blog has been updated %v", res)
}

func deleteBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Deleting blog")

	_, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: "6086ffdd4e6d91ce1488d041"})
	if err != nil {
		log.Printf("Error %v", err)
		return
	}

	fmt.Printf("Blog has been deleted")
}

func listBlogs(c blogpb.BlogServiceClient) {

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling rpc %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading %v", err)
		}

		fmt.Println(response.GetBlog())
	}
}
