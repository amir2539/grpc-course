package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	//biDiStreaming(c)
	doUnaryCallWithDeadline(c, 5*time.Second)
	doUnaryCallWithDeadline(c, 1*time.Second)
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client Started")

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Amir",
				LastName:  "asdsad",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ali",
				LastName:  "asdsad",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Asghar",
				LastName:  "asdsad",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "DFgdfg",
				LastName:  "asdsad",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error client %v", err)
	}

	for _, req := range request {
		fmt.Printf("Sending %v", req)
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	fmt.Printf("LongGreet response %v", res.GetResult())
}

func biDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Bidi  Started")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error1 %v", err)
		return
	}

	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Amir",
				LastName:  "asdsad",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ali",
				LastName:  "asdsad",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Asghar",
				LastName:  "asdsad",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "DFgdfg",
				LastName:  "asdsad",
			},
		},
	}

	// Send
	waitc := make(chan struct{})

	go func() {

		for _, req := range request {
			fmt.Printf("Sending message %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//Receive
	go func() {

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v", err)
			}

			fmt.Printf("received %v \n", res.GetResult())
		}
		close(waitc)

	}()

	<-waitc
}

func doUnaryCallWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amir",
			LastName:  "Torkaman",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {

			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("deadline exceeded")
			} else {
				log.Fatalf(" Badd Error %v", statusErr)
			}

		} else {
			log.Fatalf("Error %v", err)
		}

		return
	}

	log.Printf("Response from greet %v", res)
}
