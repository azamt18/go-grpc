package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"
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
	//makeClientStreamingCall(calculatorServiceClient)
	//makeBidirectionalStreamingCall(calculatorServiceClient)
	makeErrorUnary(calculatorServiceClient)
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

func makeBidirectionalStreamingCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting a BiDi Streaming RPC...\n")

	// create a stream by invoking the client
	stream, error := c.FindMaximum(context.Background())
	if error != nil {
		log.Fatalf("error while creating stream: %v", error)
		return
	}

	// create a wait channel to synchronize execution across goroutines
	waitChannel := make(chan struct{})

	// send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages

		// prepare a streaming request
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending message: %v \t", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}

		// close stream after sending messages
		stream.CloseSend()
	}()

	// receive a bunch of messages from the server (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while reading a stream: %v", err)
				break
			}

			fmt.Printf("Response: %v\n", res.GetMaximum())
		}

		// close channel after everything is done
		close(waitChannel)
	}()

	// block until everything is done
	<-waitChannel
}

func makeErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a SquareRoot RPC...\n")

	// correct call
	makeErrorCall(c, 10)

	// error call
	makeErrorCall(c, -2)
}

func makeErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Code())
			fmt.Printf(respErr.Message())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("Probably sent a negative number! \n")
			}
		} else {
			log.Fatalf("error while calling SquareRoot: %v \n", err)
		}
	}

	fmt.Printf("The square root of %v: %v \n", number, res.GetNumberRoot())

}
