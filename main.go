package main

import (
	"context"
	"fmt"
	"grpctls/oneway"
)

var (
	serverIp   = "127.0.0.1"
	serverPort = 5689
)

func main() {
	server := oneway.NewServer()
	server.StartServer(serverPort)
	defer server.Close()

	client := oneway.NewClient(serverIp, serverPort)
	defer client.Close()

	if person, err := client.GetPerson(context.Background(), 1); err != nil {
		panic(fmt.Sprintf("GetPersonInfo failed: %s\n", err))
	} else {
		fmt.Printf("Person: %+v\n", person)
	}
}
