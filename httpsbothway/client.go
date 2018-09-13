package httpsbothway

import (
	"net/http"
	"grpctls/common"
	"fmt"
	"io/ioutil"
)

type HttpsClient struct {
	serverAddr string
	client *http.Client
}

func NewClient(serverAddress string, commonName string) (*HttpsClient, error) {
	if commonName == "" {
		commonName = serverAddress
	}
	if config, err := common.NewClientTLSConfig(common.ClientCrt, common.ClientKey, common.CaCrt, commonName); err != nil {
		return nil, err
	}else {
		tr := &http.Transport{TLSClientConfig: config}
		return &HttpsClient{
			serverAddr: serverAddress,
			client: &http.Client{Transport: tr},
		}, nil
	}
}

func (c *HttpsClient) GetUserInfoById(id int32) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("https://%s/user/%d", c.serverAddr, id))
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()

	buf, err := ioutil.ReadAll(body)
	return string(buf), err
}
