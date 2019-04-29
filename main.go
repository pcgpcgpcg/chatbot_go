package main

import (
	pb "chatbot_go/pbx"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
)

const (
	address = "localhost:6061"
)

func SendClientHi(nmc pb.Node_MessageLoopClient){
	reqBody := new(ClientComMessage)
	clientHi:=new(MsgClientHi)
	clientHi.Id = "gRPC"
	clientHi.Version="0.15.8-rc2"
	clientHi.DeviceID="L1iC2"
	reqBody.Hi=clientHi
	//.SayHello(context.Background(), reqBody)
	err:=nmc.Send(pbCliSerialize(reqBody))
	serverMsg,err:=nmc.Recv()
	if err != nil{
		fmt.Print(err)
		return
	}
	fmt.Print(serverMsg)
}

func SendClientLogin(nmc pb.Node_MessageLoopClient){
	reqBody := new(ClientComMessage)
	clientLogin:=new(MsgClientLogin)
	clientLogin.Id = "gRPC"
	clientLogin.Scheme="basic"
	encodedString:=base64.StdEncoding.EncodeToString([]byte("alice:alice123"))
	clientLogin.Secret=[]byte(encodedString)
	reqBody.Login=clientLogin
	//.SayHello(context.Background(), reqBody)
	err:=nmc.Send(pbCliSerialize(reqBody))
	serverMsg,err:=nmc.Recv()
	if err != nil{
		fmt.Print(err)
		return
	}
	fmt.Print(serverMsg)
}

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
	nmc, err := client.MessageLoop(context.Background())
	if err != nil {
		grpclog.Fatalln(err)
	}
	// 调用方法
	reqBody := new(ClientComMessage)
	clientLogin:=new(MsgClientLogin)
	clientLogin.Id = "gRPC"
	clientLogin.Scheme="basic"
	encodedString:=base64.StdEncoding.EncodeToString([]byte("alice:alice123"))
	clientLogin.Secret=[]byte(encodedString)
	reqBody.Login=clientLogin
	//.SayHello(context.Background(), reqBody)
	err=nmc.Send(pbCliSerialize(reqBody))
	serverMsg,err:=nmc.Recv()
	fmt.Print(serverMsg)

	select {} // 阻塞
}

