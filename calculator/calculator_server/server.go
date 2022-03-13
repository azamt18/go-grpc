package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"
	"time"
)

type server struct{}

func (s server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", request)

	number := request.GetNumber() // number to be decomposed
	divisor := int64(2)           // divisor number

	for number > 1 {
		remainder := number % divisor
		if remainder == 0 {
			response := &calculatorpb.PrimeNumberDecompositionResponse{PrimeFactor: divisor}
			stream.SendMsg(response)
			time.Sleep(1000 * time.Millisecond)

			log.Printf("Factor found: %v", divisor)
			number = number / divisor
		} else {
			divisor++
			log.Printf("Divisor has increased to %v\n", divisor)
		}
	}

	return nil
}

func (s server) Sum(_ context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v\n", request)

	firstNumber := request.GetFirstNumber()
	secondNumber := request.GetSecondNumber()

	result := firstNumber + secondNumber

	response := &calculatorpb.SumResponse{Result: result}

	return response, nil
}

func main() {
	fmt.Printf("start calc server\n")

	var listener, error = net.Listen("tcp", "0.0.0.0:50052")
	if error != nil {
		log.Fatalf("failed to listen %v", error)
	}

	grpcServer := grpc.NewServer()
	calculatorServiceServer := &server{}
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, calculatorServiceServer)

	if error := grpcServer.Serve(listener); error != nil {
		log.Fatalf("failed to serve: %v\n", error)
	}
}
