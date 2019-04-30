package main

import (
	pb "chatbot_go/pbx"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
)

const (
	//address = "localhost:6061"//"192.168.8.122:6061"
	address = "192.168.8.122:6061"
)

func SendClientHi(nmc *pb.Node_MessageLoopClient) error{
	reqBody := new(ClientComMessage)
	clientHi:=new(MsgClientHi)
	clientHi.Id = "69990"
	clientHi.UserAgent="TinodeWeb/0.15.14 (Chrome/73.0; Win32); tinodejs/0.15.14"
	clientHi.Version="0.15.14"
	clientHi.DeviceID="L1iC2"
	reqBody.Hi=clientHi
	//.SayHello(context.Background(), reqBody)
	err:=(*nmc).Send(pbCliSerialize(reqBody))
	if err!=nil{
		return nil
	}

	serverMsg,err:=(*nmc).Recv()
	if err!=nil{

	}
	fmt.Println(serverMsg)
	return err
}

func SendClientLogin(nmc *pb.Node_MessageLoopClient){
	reqBody := new(ClientComMessage)
	clientLogin:=new(MsgClientLogin)
	clientLogin.Id = "69991"
	clientLogin.Scheme="basic"
	//encodedString:=base64.URLEncoding.EncodeToString([]byte("alice:alice123"))
	encodedString:="alice:alice123"
	fmt.Println(encodedString)
	clientLogin.Secret=[]byte(encodedString)
	reqBody.Login=clientLogin
	jsons,_:=json.Marshal(*clientLogin)
	fmt.Println(string(jsons))
	//.SayHello(context.Background(), reqBody)
	err:=(*nmc).Send(pbCliSerialize(reqBody))
	serverMsg,err:=(*nmc).Recv()
	if err != nil{
		fmt.Print(err)
		return
	}
	fmt.Println(serverMsg)
}

func sendClientSub(nmc *pb.Node_MessageLoopClient) error{
	reqBody := new(ClientComMessage)
	clientSub:=new(MsgClientSub)
	clientSub.Id = "69992"
	clientSub.Set(
	//encodedString:=base64.URLEncoding.EncodeToString([]byte("alice:alice123"))
	encodedString:="alice:alice123"
	fmt.Println(encodedString)
	clientLogin.Secret=[]byte(encodedString)
	reqBody.Login=clientLogin
	jsons,_:=json.Marshal(*clientLogin)
	fmt.Println(string(jsons))
	//.SayHello(context.Background(), reqBody)
	err:=(*nmc).Send(pbCliSerialize(reqBody))
	serverMsg,err:=(*nmc).Recv()
	if err != nil{
		fmt.Print(err)
		return
	}
	fmt.Println(serverMsg)
}

func main() {
	fmt.Println("starting the golang chatbot.....")
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	fmt.Println("server connected with:"+address)
	// Creates a new CustomerClient
    client:=pb.NewNodeClient(conn)
	nmc, err := client.MessageLoop(context.Background())
	if err != nil {
		grpclog.Fatalln(err)
	}
	SendClientHi(&nmc)
	SendClientLogin(&nmc)

	select {} // 阻塞
}

