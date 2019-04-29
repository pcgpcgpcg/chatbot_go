package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	pb "chatbot_go/pbx"
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
	reqBody := new(ClientComMessage)
	clientLogin:=new(MsgClientLogin)
	clientLogin.Id = "gRPC"
	clientLogin.Scheme="basic"
	clientLogin.Secret=[]byte("111")
	reqBody.Login=clientLogin


	nmc, err := client.MessageLoop(context.Background())//.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	err=nmc.Send(pbCliSerialize(reqBody))
}

