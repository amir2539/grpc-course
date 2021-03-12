package main

import (
	"context"
	"fmt"
	"github.com/amir2539/grpc/calculator/calculator_pb"
	"google.golang.org/grpc"
	"log"
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

func main() {

	fmt.Println("calculator Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculator_pb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

}
