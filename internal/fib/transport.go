package fib

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var BadNumber = errors.New("bad number in request")

func MakeHandler(ss Service) http.Handler {
	fibHandler := kithttp.NewServer(
		makeFibEndpoint(ss),
		decodeFibRequest,
		encodeResponse,
	)

	r := mux.NewRouter()

	r.Path("/fib/{n}").Handler(fibHandler).Methods("GET")

	return r
}

func decodeFibRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ns, ok := vars["n"]
	if !ok {
		return nil, BadNumber
	}

	n, err := strconv.ParseUint(ns, 10, 0)
	if err != nil {
		return nil, BadNumber
	}
	return fibRequest{N: n}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
