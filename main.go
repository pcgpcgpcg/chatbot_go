package main

import (
	pb "chatbot_go/pbx"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"strconv"
)

type UserInfo struct{
	Username string `json:"username"`
	Passwd string `json:"passwd"`
}
const (
	address = "localhost:6061"//"192.168.8.122:6061"
	//address = "192.168.8.122:6061"
)
//Hi
func sendClientHi(nmc *pb.Node_MessageLoopClient) error{
	reqBody := new(ClientComMessage)
	clientHi:=new(MsgClientHi)
	//clientHi.Id = "69990"
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
//登录
func sendClientLogin(nmc *pb.Node_MessageLoopClient,username string, passwd string ){
	reqBody := new(ClientComMessage)
	clientLogin:=new(MsgClientLogin)
	//clientLogin.Id = "69991"
	clientLogin.Scheme="basic"
	//encodedString:=base64.URLEncoding.EncodeToString([]byte("alice:alice123"))
	//encodedString:="alice:alice123"
	encodedString:=username+":"+passwd
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
//进入编组
func sendClientSub(nmc *pb.Node_MessageLoopClient,topic_name string) error{
	reqBody := new(ClientComMessage)
	clientSub:=new(MsgClientSub)
	clientSub.Topic=topic_name
	//add set
	clientMsgSetQuery:=new(MsgSetQuery)
	clientMsgSetSub:=new(MsgSetSub)
	clientMsgSetSub.Mode="JRWPS"
	clientMsgSetQuery.Sub=clientMsgSetSub
	clientSub.Set=clientMsgSetQuery
	//add sub
	clientMsgGetQuery:=new(MsgGetQuery)
	clientMsgGetOpts:=new(MsgGetOpts)
	clientMsgGetOpts.Limit=24
	clientMsgGetQuery.Data=clientMsgGetOpts
	clientMsgGetQuery.What="data sub desc"

	clientSub.Set=clientMsgSetQuery
	clientSub.Get=clientMsgGetQuery

	reqBody.Sub=clientSub
	//send
	err:=(*nmc).Send(pbCliSerialize(reqBody))
	//recv
	serverMsg,err:=(*nmc).Recv()
	if err != nil{
		fmt.Print(err)
		return err
	}
	fmt.Println(serverMsg)
	return nil
}
//发送消息
func sendClientPub(nmc *pb.Node_MessageLoopClient,topic_name string,msg string) error{
	reqBody := new(ClientComMessage)
	clientPub:=new(MsgClientPub)
	clientPub.Topic=topic_name
	clientPub.NoEcho=true
	clientPub.Content=msg

	reqBody.Pub=clientPub
	//send
	err:=(*nmc).Send(pbCliSerialize(reqBody))
	//recv
	serverMsg,err:=(*nmc).Recv()
	if err != nil{
		fmt.Print(err)
		return err
	}
	fmt.Println(serverMsg)
	return nil
}

func main() {
	/*fmt.Println("starting the golang chatbot.....")
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
	//初始化交互
	sendClientHi(&nmc)
	//登录
	sendClientLogin(&nmc,"alice","alice123")
	//切入编组
	sendClientSub(&nmc,"grp_3cPvkpRTy8")
	//发送消息
	sendClientPub(&nmc,"grp_3cPvkpRTy8","Hi,I am Bot1.")*/

	testManyConnection()

	select {} // 阻塞
}

func loginToServer(username string,passwd string,index int){
	fmt.Println("starting the golang chatbot "+strconv.Itoa(index))
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
	sendClientHi(&nmc)
	sendClientLogin(&nmc,username,passwd)
	sendClientSub(&nmc,"grpQOnSKtB1iPo")
	sendClientPub(&nmc,"grpQOnSKtB1iPo","Hi,I am Bot "+strconv.Itoa(index))
}


func testManyConnection(){
	var users=make([]UserInfo, 200)
	for i:=0;i<200;i++{
		users[i]=UserInfo{"bot"+strconv.Itoa(i),"123456"}
	}

	for index,value:=range users{
		loginToServer(value.Username,value.Passwd,index)
	}
}

