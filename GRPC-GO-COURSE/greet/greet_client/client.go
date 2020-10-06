package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vamsikrishnasiddu/gRPC-go-Projects/GRPC-GO-COURSE/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect:%v", err)
	}

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBidiStreaming(c)
	//fmt.Printf("Created client %f", c)
	doUnaryWithDeadline(c, 5*time.Second) //should complete
	doUnaryWithDeadline(c, 1*time.Second) //should timeout
}
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "stephane",
			LastName:  "Marek",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Println(res.Result)

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "stephane",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Joseph",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet:%v", err)

	}
	//we iterate over our slice and send each message individually.
	for _, req := range requests {
		fmt.Printf("Sending req:%v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while reciveing response from LongGreet:%v", err)
	}

	fmt.Printf("LongGreet Response:%v\n", res)

}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bidirectional RPC....")
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "stephane",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Joseph",
			},
		},
	}
	//we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream:%v", err)
		return
	}
	//we send a bunch of messages to the client(go routine)
	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Println("Sending message:", req)
			err := stream.Send(req)
			time.Sleep(1000 * time.Millisecond)

			if err != nil {
				log.Fatalf("Error while sending the data to the stream%v\n", err)
			}
		}
		stream.CloseSend()

	}()

	//we receive a bunch of messages from the client(go routine)

	go func() {
		//function to recieve a bunch of messages
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break

			}

			if err != nil {
				log.Fatalf("Error while recieving data from server.%v\n", err)
			}

			fmt.Printf("Recieved;%v\n", res)

		}
		close(waitc)

	}()

	//block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary With Deadline RPC.. call")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {

		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit:! Deadline exceeded.")
			} else {
				fmt.Println("Unexpected Error:", err)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline Rpc:%v", err)
		}
		return
	}

	log.Printf("Response from GreetWithDeadline:%v", res.Result)
}
