package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Client Started")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connet %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amir",
			LastName:  "Torkaman",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	log.Printf("Response from greet %v", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client many")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amir",
			LastName:  "Torkaman",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// Server closed stream
			break
		}

		if err != nil {
			log.Fatalf("Error Stream %v", err)
		}

		log.Printf("Result from manytimes %v", msg.GetResult())
	}
}
