package main

import (
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

// type ImplementCalculator struct {
// }

// func (I *ImplementCalculator) Sum(ctx context.Context, in *calculatorpb.CalculatorRequest, opts ...grpc.CallOption) (*calculatorpb.CalculatorResponse, error) {
// 	// return nil, status.Errorf(codes.Unimplemented, "method Greet not implemented")
// 	fmt.Printf("Greet function was invoked with %v\n", in)
// 	res := &calculatorpb.CalculatorResponse{
// 		Result: 123,
// 	}
// 	return res, nil
// }

func main() {
	fmt.Println("Hello World this is calculator server")

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &calculatorpb.UnimplementedCalculatorServiceServer{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
