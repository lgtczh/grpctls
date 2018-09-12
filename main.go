package main

import (
	"context"
	"fmt"
	"grpctls/protos"
	"grpctls/replica"
)

var (
	serverIp   = "127.0.0.1"
	serverPort = 5689
)

func main() {
	server := replica.New()
	server.StartServer(serverPort)
	defer server.Close()

	conn := replica.NewClientConn(serverIp, serverPort)
	defer conn.Close()

	client := protos.NewPersonInfoProviderClient(conn)
	if person, err := client.GetPersonInfo(context.Background(), &protos.Request{ReqId: 1}); err != nil {
		panic(fmt.Sprintf("GetPersonInfo failed: %s\n", err))
	} else {
		fmt.Printf("Person: %+v\n", person)
	}
}
