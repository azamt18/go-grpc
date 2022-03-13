package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
)

func main() {
	fmt.Printf("client started")

	var connection, error = grpc.Dial("localhost:50052", grpc.WithInsecure())
	if error != nil {
		log.Fatalf("could not connect: %v", error)
	}

	defer func(connection *grpc.ClientConn) {
		err := connection.Close()
		if err != nil {
			log.Fatalf("error while closing connection: %v", err)
		}
	}(connection)

	calcServiceClient := calculatorpb.NewCalculatorServiceClient(connection)

	makeUnaryCall(calcServiceClient)

}

func makeUnaryCall(calcServiceClient calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do an unary RPC \n")

	// prepare request
	request := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 5,
	}

	// make a request
	var response, error = calcServiceClient.Sum(context.Background(), request)

	// handle a request
	if error != nil {
		log.Fatalf("error while calling Calc RPC: %v", error)
	}

	log.Printf("response from Calc RPC: %v", response.Result)
}
