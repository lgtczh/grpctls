package grpconeway

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpctls/common"
	"grpctls/protos"
)

type Client struct {
	clientConn *grpc.ClientConn
	grpcClient protos.PersonInfoProviderClient
}

func newClientConn(serverIp string, serverPort int, commonName string) *grpc.ClientConn {
	var opts []grpc.DialOption
	if commonName == "" {
		commonName = serverIp
	}

	cred, err := credentials.NewClientTLSFromFile(common.CaCrt, commonName)
	if err != nil {
		panic(fmt.Sprintf("Client: failed to create TLS credentials %v", err))
	}

	opts = append(opts, grpc.WithTransportCredentials(cred))
	masterAddr := fmt.Sprintf("%s:%d", serverIp, serverPort)

	conn, err := grpc.Dial(masterAddr, opts...)
	if err != nil {
		panic(fmt.Sprintf("Client: can't connect to %s: %v", masterAddr, err))
	}

	return conn
}

func NewClient(serverIp string, serverPort int) *Client {
	conn := newClientConn(serverIp, serverPort, "")
	return &Client{
		clientConn: conn,
		grpcClient: protos.NewPersonInfoProviderClient(conn),
	}
}

func (c *Client) Close() {
	c.clientConn.Close()
}

func (c *Client) GetPerson(ctx context.Context, id int32) (*protos.Person, error) {
	return c.grpcClient.GetPersonInfo(ctx, &protos.Request{ReqId: id})
}
