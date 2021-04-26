package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/calculator/calculator_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculator_pb.SumRequest) (*calculator_pb.SumResponse, error) {
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber
	res := &calculator_pb.SumResponse{SumResult: sum}
	return res, nil
}
func (*server) PrimeNumberDeComp(req *calculator_pb.PrimeNumberDeCompRequest, stream calculator_pb.CalculatorService_PrimeNumberDeCompServer) error {
	fmt.Print("recdasd;")
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculator_pb.PrimeNumberDeCompResponse{
				PrimeFactor: divisor,
			})
			number /= divisor
		} else {
			divisor++
			fmt.Printf("divisor increased %v", divisor)
		}
	}
	return nil
}
func (*server) ComputeAverage(stream calculator_pb.CalculatorService_ComputeAverageServer) error {

	fmt.Println("recieved average server\n")

	sum := int64(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculator_pb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error occiured %v", err)
		}

		sum += req.GetNumber()
		count++
	}
}
func (*server) FindMaximum(stream calculator_pb.CalculatorService_FindMaximumServer) error {

	fmt.Printf("Find max server")
	maximumNubmer := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error when reading %v", err)
		}

		number := req.GetNumber()
		if number > maximumNubmer {
			maximumNubmer = number
			err := stream.Send(&calculator_pb.FindMaximumResponse{Maximum: maximumNubmer})
			if err != nil {
				log.Fatalf("Error when sending data %v", err)

			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculator_pb.SquareRootRequest) (*calculator_pb.SquareRootResponse, error) {
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number: %v", number),
		)
	}

	return &calculator_pb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {

	fmt.Println("calculator Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculator_pb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

}
