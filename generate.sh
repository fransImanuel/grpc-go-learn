protoc .\greet\greetpb\greet.proto --go-grpc_out=.\ .\greet\greetpb\greet.proto
protoc .\calculator\calculatorpb\calculator.proto --go-grpc_out=.\ .\calculator\calculatorpb\calculator.proto

protoc greet\greetpb\greet.proto --go-grpc_out=.\ .\greet\greetpb\greet.proto --go_out=.\ greet\greetpb\greet.proto
protoc .\calculator\calculatorpb\calculator.proto --go_out=.\ .\calculator\calculatorpb\calculator.proto --go-grpc_out=.\ .\calculator\calculatorpb\calculator.proto