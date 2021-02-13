package ip

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(is Service) http.Handler {
	ipHandler := kithttp.NewServer(
		makeIpEndpoint(is),
		decodeIpRequest,
		encodeResponse,
	)

	r := mux.NewRouter()

	r.Path("/ip").Handler(ipHandler).Methods("GET")

	return r
}

func decodeIpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
