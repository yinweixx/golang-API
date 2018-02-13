package api

import(
	"fmt"
	"log"
	"context"
	"google.golang.org/grpc"
	"e.coding.net/anyun-cloud-api-gateway/server"
	pb "e.coding.net/anyun-cloud-api-gateway/apigrpc"
)

const(
	address = "localhost:50051"
)

func sendHelp(api *server.APICONTROLLERPARAMS){
	fmt.Println(api.ID)
	conn,err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect:%v",err)
	}
	defer conn.Close()
	c := pb.NewSendHelpClient(conn)

	r,err := c.SendHelp(context.Background(),&pb.Request{
		Id		: api.ID,
		Name	: api.Name,
		Version : api.Version,
		Dc		: api.Dc,
	})
	if err != nil{
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Data)
	log.Printf("Greeting: %s", r.Id)
}