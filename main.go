package main

import (
	"context"
	"fmt"
	"grpctls/grpconeway"
	"grpctls/httpsbothway"
)

var (
	serverIp       = "127.0.0.1"
	grpcServerPort = 5689
	httpsServerPort = 8000
)

func main() {
	server := grpconeway.NewServer()
	server.StartServer(grpcServerPort)
	defer server.Close()

	client := grpconeway.NewClient(serverIp, grpcServerPort)
	defer client.Close()

	if person, err := client.GetPerson(context.Background(), 1); err != nil {
		panic(fmt.Sprintf("GetPersonInfo failed: %s\n", err))
	} else {
		fmt.Printf("Person: %+v\n", person)
	}

	httpsServer := httpsbothway.NewServer(httpsServerPort)
	httpsServer.StartServer()
	defer httpsServer.Close()

	httpsClient, err := httpsbothway.NewClient(fmt.Sprintf("%s:%d", serverIp, httpsServerPort), serverIp)
	if err != nil {
		panic(err)
	}
	if info, err := httpsClient.GetUserInfoById(1); err != nil {
		panic(err)
	} else {
		fmt.Printf("Person: %s\n", info)
	}
}
