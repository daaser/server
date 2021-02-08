package header

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeHeaderEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*http.Request)
		headers := svc.Headers(req)
		return headers, nil
	}
}
