package header

import (
	"bytes"
	"net/http"
	"strings"
)

type Service interface {
	Headers(*http.Request) []byte
}

type service struct{}

func (svc service) Headers(req *http.Request) []byte {
	var buf bytes.Buffer
	for hk, hv := range req.Header {
		buf.WriteString(hk)
		buf.WriteString(": ")
		buf.WriteString(strings.Join(hv, ","))
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

func NewService() Service {
	return &service{}
}
