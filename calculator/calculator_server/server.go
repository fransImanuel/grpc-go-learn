package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (server) Sum(ctx context.Context, in *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
	fmt.Printf("Received Sum RPC:%v\n", in)
	num1 := in.GetNum1()
	num2 := in.GetNum2()
	result := num1 + num2
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}
func (server) PrimeNumber(req *calculatorpb.CalculatorPrimeRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	// return status.Errorf(codes.Unimplemented, "method PrimeNumber not implemented")
	fmt.Printf("Received Sum RPC:%v\n", req)

	//logic for find prime decomposition
	var k int64 = 2
	var number int64 = req.Number
	for number > 1 {
		if number%k == 0 {
			// fmt.Println(k)
			res := &calculatorpb.CalculatorPrimeResponse{Number: k}
			stream.Send(res)
			number = number / k
		} else {
			k++
		}
	}
	return nil
}
func (server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with a streaming request\n")
	var result float64
	var i float64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result = result / i
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				AVGResult: int64(result),
			})
		}
		if err != nil {
			log.Fatal("Error while reading client stream: %v", err)
		}
		i++
		number := req.GetNumber()
		result = result + float64(number)
	}
}
func (server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a Bidi request\n")
	var numberStub []int
	var temp int = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		inputNumber := req.GetNumber()

		//check input (is this input bigger than before or not)
		numberStub = append(numberStub, int(inputNumber))
		for _, number := range numberStub {
			if temp <= number {
				temp = number
			}
		}

		if temp <= int(inputNumber) {
			stream.Send(&calculatorpb.FindMaxResponse{
				Number: inputNumber,
			})
		}
	}
}
func (server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method SquareRoot not implemented")
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recevied a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: float32(math.Sqrt(float64(number))),
	}, nil
}

func main() {
	fmt.Println("Hello World this is calculator server")

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
