package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	pb "chatbot_go/proto"
)

const (
	address = "localhost:6061"
)

func main() {
	fmt.Print("starting the golang chatbot.....")
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	fmt.Print("server connected with:"+address)
	// Creates a new CustomerClient
client:=pb.NewNodeClient(conn)
	// 调用方法
	reqBody := new(pb.ClientLogin)
	reqBody.Id = "gRPC"
	reqBody.Scheme="basic"
	reqBody.Secret=[]byte("111")
	nmc, err := client.MessageLoop(context.Background())//.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	err=nmc.Send(&pb.ClientMsg{login:{Id:"gRPC",Scheme:"basic"})
}

