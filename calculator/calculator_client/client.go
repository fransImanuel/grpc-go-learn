package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm calculator client")

	cc, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created client: %f", c)
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBidiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Num1: 3,
		Num2: 10,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}
	log.Printf("Result = %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &calculatorpb.CalculatorPrimeRequest{Number: 120}
	resStream, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//when reached the end of the stream
			break
		}
		if err != nil {
			log.Fatal("error while reading stream: %v", err)
		}
		log.Printf("Response from Greet: %v", msg.Number)
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	request := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatal("error while calling ComputeAverage: %v", err)
	}

	//we iterate over our slice and send each message individually
	for _, req := range request {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bidi Streaming RPC...")

	//we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating streaming: %v", err)
	}

	numbers := []int32{4, 7, 2, 19, 4, 6, 32}

	waitc := make(chan struct{})
	//we send a bunch of messages to the client (go routine)
	go func() {
		//function to send bunch of messsages
		for _, number := range numbers {
			fmt.Printf("Sendding Message: %v\n", number)
			stream.Send(&calculatorpb.FindMaxRequest{
				Number: int64(number),
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		for {
			//function to receive a bunch of messages
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receive: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetNumber())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	//correct call
	doErrorCall(c, 10)

	//error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: int64(n)})

	if err != nil {
		respError, ok := status.FromError(err)
		if ok {
			fmt.Printf("Error Message From Server : %v\n", respError.Message())
			fmt.Println(respError.Code())
			if respError.Code() == codes.InvalidArgument {
				fmt.Println("We Probably Sent a Negative Number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling squareroot: %v", err)
			return
		}
	}
	fmt.Printf("Result of SquareRoot of %v : %v\n", n, res.GetNumberRoot())
}
