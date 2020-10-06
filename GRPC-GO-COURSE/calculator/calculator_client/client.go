package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vamsikrishnasiddu/gRPC-go-Projects/GRPC-GO-COURSE/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Cannot Establish connectin :%v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBidiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		Number1: 1,
		Number2: 2,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("The sum of two numbers is :", res.Result)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.PrimeRequest{
		Number: 120,
	}

	resStream, err := c.Prime(context.Background(), req)

	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("The Response from Prime RPC server:", msg.Number)

	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Numberavg: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())
}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximum Api with BidiStreaming RPC...")

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 5,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 20,
		},
	}

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error creating the client Stream:%v\n", err)
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending number:%v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)

		}

		stream.CloseSend()

	}()

	go func() {

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break

			}

			if err != nil {
				log.Fatalf("error while recieving the data from the server:%v\n", err)
				break
			}
			fmt.Println("The Response from the Server is:", res.Maximum)
		}
		close(waitc)

	}()

	<-waitc

}
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	doErrorCall(c, 10)
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)

		if ok {
			//actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())

			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling squareRoot:%v\n", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v:%v\n", n, res.GetNumberRoot())
}
