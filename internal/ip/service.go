package ip

import (
	"io/ioutil"
	"net/http"
)

type Service interface {
	GetIp() ([]byte, error)
}

type service struct {
	client *http.Client
}

func (svc service) GetIp() ([]byte, error) {
	resp, err := svc.client.Get("https://api.ipify.org")
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(resp.Body)
}

func NewService() Service {
	return &service{client: &http.Client{}}
}
