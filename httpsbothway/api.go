package httpsbothway

import (
	"grpctls/common"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"grpctls/protos"
	"fmt"
	"encoding/json"
)
type UserAPI struct {
	routes common.Routes
	server *HttpsServer
}

type Response struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func NewAPI(srv *HttpsServer) *UserAPI{
	u := &UserAPI{
		server: srv,
	}
	u.routes = common.Routes{
			common.Route{
				Name:        "UserInfo",
				Method:      "GET",
				Pattern:     "/user/{id}",
				HandlerFunc: u.AllAgentGet,
			},
		}
	return u
}

func (u *UserAPI) AllAgentGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var statusCode int
	var user *protos.Person
	var err error
	if id, ok := mux.Vars(r)["id"]; ok {
		if v, e := strconv.Atoi(id); e == nil {
			if user, e = u.server.GetUserById(int32(v)); e == nil {
				statusCode = 200
			} else {
				err = e
				statusCode = 404
			}
		} else {
			err = e
			statusCode = 400
		}
	} else {
		err = fmt.Errorf("Id must give\n")
		statusCode = 400
	}
	resp := Response{Status:statusCode, Data:user}
	if err != nil {
		resp.Error = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}

func (u *UserAPI) NewAPIRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range u.routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}