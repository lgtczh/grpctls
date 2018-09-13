package httpsbothway

import (
	"grpctls/common"
	"net/http"
	"fmt"
	"grpctls/protos"
)

type HttpsServer struct {
	users []*common.User
	apiServer *http.Server
}

func NewServer(port int) *HttpsServer {
	server := &HttpsServer{users: common.Users}
	handler := NewAPI(server).NewAPIRouter()
	server.apiServer = &http.Server{
		Handler: handler,
		Addr: fmt.Sprintf(":%d", port),
	}
	return server
}

func (s *HttpsServer) GetUserById(id int32) (*protos.Person, error){
	for _, user := range s.users {
		if user.GetId() == id {
			return &protos.Person{Id: user.GetId(), Name: user.GetName(), Email: user.GetEmail()}, nil
		}
	}
	return nil, fmt.Errorf("Not found %d\n", id)
}

func (s *HttpsServer) StartServer() error {
	config, err := common.NewServerTlsConfig(common.CaCrt)
	if err != nil {
		return err
	}
	s.apiServer.TLSConfig = config
	go func() {
		if err := s.apiServer.ListenAndServeTLS(common.ServerCrt, common.ServerKey); err != http.ErrServerClosed {
			panic("Unexpected error: " + err.Error())
		}
	}()
	return nil
}

func (s *HttpsServer) Close() {
	s.apiServer.Close()
}
