package replica

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpctls/protos"
	"net"
)

type Server struct {
	personInfo []*User
	crt        string
	key        string
	grpcServer *grpc.Server
}

type User struct {
	id    int32
	name  string
	email string
}

func New() *Server {
	users := []*User{
		{1, "Tom", "tom@email.com"},
		{2, "Jack", "jack@email.com"},
	}
	server := &Server{
		personInfo: users,
		crt:        serverCrt,
		key:        serverKey,
	}
	creds := server.getServerCreds()
	server.grpcServer = grpc.NewServer(grpc.Creds(creds))

	return server
}

func (s *Server) StartServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(fmt.Sprintf("Server: failed to listen: %v", err))
	}

	protos.RegisterPersonInfoProviderServer(s.grpcServer, s)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			panic(fmt.Sprintf("grpcServer error: %s\n", err))
		}
	}()
}

func (s *Server) Close() {
	s.grpcServer.Stop()
}

func (s *Server) getServerCreds() credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(s.crt, s.key)
	if err != nil {
		panic(fmt.Sprintf("Server: failed to generate credentials: %s\n", err))
	}

	return creds
}

func (s *Server) GetPersonInfo(ctx context.Context, in *protos.Request) (*protos.Person, error) {
	for _, user := range s.personInfo {
		if user.id == in.ReqId {
			return &protos.Person{Id: user.id, Name: user.name, Email: user.email}, nil
		}
	}
	return nil, fmt.Errorf("Not found %d\n", in.ReqId)
}
