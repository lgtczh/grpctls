package common

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"net/http"
)

var (
	ServerCrt = "tls/server/127.0.0.1.crt"
	ServerKey = "tls/server/127.0.0.1.key"
	ClientCrt = "tls/client/127.0.0.1.crt"
	ClientKey = "tls/client/127.0.0.1.key"
	CaCrt = "tls/ca/ca.crt"
	Users = []*User{
		{1, "Tom", "tom@email.com"},
		{2, "Jack", "jack@email.com"},
	}
)

type User struct {
	id    int32
	name  string
	email string
}

func (u *User) GetId() int32 {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewServerTlsConfig(caCrt string) (*tls.Config, error) {
	certPool, err := newCACertPool(caCrt)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientAuth:             tls.RequireAndVerifyClientCert,
		ClientCAs:              certPool,
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
	}, nil
}

func newCACertPool(caCrt string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCrt)
	if err != nil {
		return nil, fmt.Errorf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("failed to append ca certs %s", caCrt)
	}
	return certPool, nil
}

func NewClientTLSConfig(clientCRT, clientKey, rootCRT, serverName string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(clientCRT, clientKey)
	if err != nil {
		return nil, fmt.Errorf("could not read client certificate: %s", err)
	}
	certPool, err := newCACertPool(rootCRT)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}, nil
}

