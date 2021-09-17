package main

import (
	"fmt"
	"log"
	"net"

	"grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

// type server struct {
// }

// func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest, opts ...grpc.CallOption) (*greetpb.GreetResponse, error) {
// 	fmt.Printf("Greet function was invoked with %v\n", in)
// 	firstName := in.GetGreeting().GetFirstName()
// 	result := "Hello " + firstName
// 	res := &greetpb.GreetResponse{
// 		Result: result,
// 	}
// 	return res, nil
// }

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, greetpb.UnimplementedGreetServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
