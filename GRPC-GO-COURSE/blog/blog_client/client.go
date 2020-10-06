package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/vamsikrishnasiddu/gRPC-go-Projects/GRPC-GO-COURSE/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client...")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	defer cc.Close()
	if err != nil {
		log.Fatalf("Couldn't connect:%v", err)
	}

	c := blogpb.NewBlogServiceClient(cc)
	//Create blog
	fmt.Println("Creating the blog...")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "my first blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Unexpected Err:%v\n", err)
		return
	}
	fmt.Printf("Blog has been Created %v\n", createBlogRes)

	blogID := createBlogRes.GetBlog().GetId()

	fmt.Println("Reading the blog..")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "dfsfkfsfsflsf"})

	if err2 != nil {
		fmt.Printf("Error while reading:%v\n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)

	if readBlogErr != nil {
		fmt.Printf("Error while reading:%v\n", readBlogErr)
	}

	fmt.Printf("Blog was Read:%v\n", readBlogRes)

	//update the blog

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Stephane Updated",
		Title:    "Changed Author",
		Content:  "Content of the first blog Added Awesome features",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if updateErr != nil {
		log.Fatalf("Error while calling the update method:%v\n", updateErr)

	}

	fmt.Printf("Blog was updated:%v\n", updateRes)

	//deleting the blog

	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		log.Fatalf("Error while deleting:%v\n", deleteErr)
	}

	fmt.Printf("Blog was deleted:%v\n", deleteRes)

	//list Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v\n", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened:%v", err)
		}

		fmt.Println(res.GetBlog())
	}
}
