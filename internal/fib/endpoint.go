package fib

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type fibRequest struct {
	N uint64
}

type fibResponse struct {
	Fib string `json:"fib,omitempty"`
}

func makeFibEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(fibRequest)
		fib := svc.Fib(req.N)
		// return fibResponse{fib.String()}, nil
		return fib, nil
	}
}
