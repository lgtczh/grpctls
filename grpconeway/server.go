package grpconeway

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpctls/common"
	"grpctls/protos"
	"net"
)

type Server struct {
	personInfo []*common.User
	grpcServer *grpc.Server
}

func NewServer() *Server {

	return &Server{
		personInfo: common.Users,
		grpcServer: grpc.NewServer(grpc.Creds(getServerCred(common.ServerCrt, common.ServerKey))),
	}
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

func getServerCred(crtFile, keyFile string) credentials.TransportCredentials {
	cred, err := credentials.NewServerTLSFromFile(crtFile, keyFile)
	if err != nil {
		panic(fmt.Sprintf("Server: failed to generate credentials: %s\n", err))
	}

	return cred
}

func (s *Server) GetPersonInfo(ctx context.Context, in *protos.Request) (*protos.Person, error) {
	for _, user := range s.personInfo {
		if user.GetId() == in.ReqId {
			return &protos.Person{Id: user.GetId(), Name: user.GetName(), Email: user.GetEmail()}, nil
		}
	}
	return nil, fmt.Errorf("Not found %d\n", in.ReqId)
}
