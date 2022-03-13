package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

func (s server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("Received FindMaximum streaming RPC\n")
	numbers := make([]int32, 0)
	max := int32(0) // smallest integer&non-negative number

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		inputNumber := req.GetNumber()
		numbers = append(numbers, inputNumber)
		if inputNumber > max {
			max = inputNumber
		}

		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: max,
		})
		if sendErr != nil {
			log.Fatalf("error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func (s server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Received ComputeAverage streaming RPC\n")
	sum := int64(0)
	length := int64(0)
	average := 0.0

	for {
		request, error := stream.Recv()
		if error == io.EOF {
			// finished reading the client stream
			average = float64(sum) / float64(length)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}

		// handle error
		if error != nil {
			log.Fatalf("error while reading client stream: %v\n", error)
		}

		sum += request.GetNumber()
		length = length + 1
	}
}

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
	fmt.Printf("Staring CalculatorService\n")

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
