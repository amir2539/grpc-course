package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/calculator/calculator_pb"
	"google.golang.org/grpc"
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
	doClientStreaming(c)
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
