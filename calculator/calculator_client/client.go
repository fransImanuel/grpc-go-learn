package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
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
	doBidiStreaming(c)
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

func doBidiStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Println("Starting to do a Bidi Streaming RPC...")

	//we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating streaming: %v",err)
	}

	requests:= []*calculatorpb.FindMaxRequest{
		&calculatorpb.FindMaxRequest{
			Number: 1,
		},
		&calculatorpb.FindMaxRequest{
			Number: 5,
		},
		&calculatorpb.FindMaxRequest{
			Number: 3,
		},
		&calculatorpb.FindMaxRequest{
			Number: 6,
		},
		&calculatorpb.FindMaxRequest{
			Number: 2,
		},
		&calculatorpb.FindMaxRequest{
			Number: 20,
		},
	}

	waitc := make(chan struct{})
	//we send a bunch of messages to the client (go routine)
	go func() {
		//function to send bunch of messsages
		for _, req := range requests {
			fmt.Printf("Sendding Message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		for{
			//function to receive a bunch of messages
			res,err := stream.Recv()
			if err == io.EOF {
				break;
			}
			if err != nil {
				log.Fatalf("Error while receive: %v",err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetNumber())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}
