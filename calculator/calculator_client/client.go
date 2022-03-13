package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
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

	calculatorServiceClient := calculatorpb.NewCalculatorServiceClient(connection)

	//makeUnaryCall(calculatorServiceClient)
	//makeServerStreamingCall(calculatorServiceClient)
	makeClientStreamingCall(calculatorServiceClient)
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

	// handle an error
	if error != nil {
		log.Fatalf("error while calling unary RPC: %v", error)
	}

	log.Printf("response from unary RPC: %v", response.Result)

	log.Printf("\t End of unary call \n\n")
}

func makeServerStreamingCall(calcServiceClient calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a PrimeDecomposition Server Streaming RPC... \n")

	// prepare request
	request := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 150,
	}

	// make a stream request
	responseStream, error := calcServiceClient.PrimeNumberDecomposition(context.Background(), request)

	// handle an error
	if error != nil {
		log.Fatalf("error while calling server streaming RPC: %v", error)
	}

	for true {
		message, error := responseStream.Recv()
		if error == io.EOF {
			// reached end of stream
			break
		}

		// handle stream error
		if error != nil {
			log.Printf("error while reading stream: %v", error)
		}

		// handle stream response
		log.Printf("PrimeFactor: %v\n", message.GetPrimeFactor())
	}

	fmt.Printf("Closing PrimeDecomposition Server Streaming RPC...\n")
}

func makeClientStreamingCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a ComputeAverage client streaming RPC... \n")

	stream, error := c.ComputeAverage(context.Background())
	numbers := []int64{1, 2, 3, 4}

	// iterate over the request and send each message individually
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	if error != nil {
		log.Fatalf("error while calling ComputeAverage: %v", error)
	}

	response, error := stream.CloseAndRecv()
	if error != nil {
		log.Fatalf("error while reading response from ComputeAverage: %v", error)
	}

	log.Printf("ComputeAverage response: %v", response.GetAverage())
}
