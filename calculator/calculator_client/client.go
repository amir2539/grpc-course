package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/calculator/calculator_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

func main() {
	fmt.Println("cal Client Started")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connet %v", err)
	}

	defer cc.Close()
	c := calculator_pb.NewCalculatorServiceClient(cc)

	//doUnary(c)
	//doServerStreamin(c)
	//doClientStreaming(c)
	//doBidiStreaming(c)
	squareRootUnary(c)
}

func doUnary(c calculator_pb.CalculatorServiceClient) {
	req := &calculator_pb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	log.Printf("Response from greet %v", res.SumResult)
}

func doServerStreamin(c calculator_pb.CalculatorServiceClient) {
	fmt.Println("asdasd")
	req := &calculator_pb.PrimeNumberDeCompRequest{
		Number: 100,
	}

	stream, err := c.PrimeNumberDeComp(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("sadasd %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(client calculator_pb.CalculatorServiceClient) {
	fmt.Println("start client average")

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error %v", err)
	}

	numbers := []int64{2, 2, 4}

	for _, number := range numbers {

		stream.Send(&calculator_pb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error2 %v", err)
	}

	fmt.Printf("average: %v \n", res.GetAverage())
}

func doBidiStreaming(client calculator_pb.CalculatorServiceClient) {
	fmt.Printf("client started")

	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error when opening strea %v", err)
	}

	waitc := make(chan struct{})

	// Send
	go func() {

		numbers := []int32{4, 656, 7, 344, 7666, 20, 4534534}
		for _, number := range numbers {
			stream.Send(&calculator_pb.FindMaximumRequest{Number: number})
			//time.Sleep(time.Millisecond * 1000)
		}
		stream.CloseSend()
	}()

	// Receive
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem when reading client %v", err)
			}

			maximum := res.GetMaximum()
			fmt.Printf("Received new max :%v \n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func squareRootUnary(c calculator_pb.CalculatorServiceClient) {

	doErrorCall(c, int32(10))

	doErrorCall(c, int32(-2))
}

func doErrorCall(c calculator_pb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculator_pb.SquareRootRequest{
		Number: number,
	})
	if err != nil {
		resErr, ok := status.FromError(err)

		if ok {
			// Is error form grpc
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())

			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("negative number")
				return
			}

		} else {
			log.Fatalf("Big Error %v", err)
		}
	}

	fmt.Printf("Result %v sqrt= %v", number, res.GetNumberRoot())
}
