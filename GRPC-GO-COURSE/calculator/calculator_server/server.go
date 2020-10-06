package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/vamsikrishnasiddu/gRPC-go-Projects/GRPC-GO-COURSE/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	num1 := req.GetNumber1()
	num2 := req.GetNumber2()

	result := num1 + num2

	res := &calculatorpb.SumResponse{
		Result: result,
	}

	return res, nil

}

func (*server) Prime(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_PrimeServer) error {

	fmt.Println("Prime function in Server is invoked with req:", req)

	num := req.GetNumber()

	var k int32 = 2
	for num > 1 {
		if num%k == 0 {
			err := stream.Send(
				&calculatorpb.PrimeResponse{
					Number: k,
				},
			)
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				log.Fatalln(err)
			}
			num = num / k
		} else {
			k = k + 1
		}
	}
	return nil

}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			//we have finished reading the client.
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})

		}
		if err != nil {
			log.Fatalf("Error in the compute Average:%v\n", err)
		}

		sum += sum + req.GetNumberavg()
		count++
	}

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum is invoked with BiDirectional streaming")
	var max int32 = math.MinInt32
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return err
		}

		if err != nil {
			log.Fatalf("Error reading the stream of data from client:%v\n", err)
			return err
		}

		if req.GetNumber() > max {
			max = req.GetNumber()
		}

		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: max,
		})

		if sendErr != nil {
			log.Fatalf("Error while sending the data to the client%v\n", sendErr)
			return sendErr
		}

	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number:%v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Server started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Couldn't start the server %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve%v\n", err)
	}
}
