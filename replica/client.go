package replica

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewClientConn(serverIp string, serverPort int) *grpc.ClientConn {
	var opts []grpc.DialOption

	creds, err := credentials.NewClientTLSFromFile(caCrt, "127.0.0.1")
	if err != nil {
		panic(fmt.Sprintf("Client: failed to create TLS credentials %v", err))
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))
	masterAddr := fmt.Sprintf("%s:%d", serverIp, serverPort)

	conn, err := grpc.Dial(masterAddr, opts...)
	if err != nil {
		panic(fmt.Sprintf("Client: can't connect to %s: %v", masterAddr, err))
	}

	return conn
}
