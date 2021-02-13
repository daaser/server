package ip

import (
	"context"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type ipResponse struct {
	IP  string `json:"ip"`
	Err string `json:"error,omitempty"`
}

func makeIpEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var builder strings.Builder
		ip, err := svc.GetIp()
		if err != nil {
			return ipResponse{"", err.Error()}, nil
		}

		_, err = builder.Write(ip)
		if err != nil {
			return ipResponse{"", err.Error()}, nil
		}

		return ipResponse{builder.String(), ""}, nil
	}
}
