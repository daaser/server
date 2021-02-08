package str

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(ss Service) http.Handler {
	uppercaseHandler := kithttp.NewServer(
		makeUppercaseEndpoint(ss),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := kithttp.NewServer(
		makeCountEndpoint(ss),
		decodeCountRequest,
		encodeResponse,
	)

	r := mux.NewRouter()

	r.Methods("POST").Path("/string/uppercase").Handler(uppercaseHandler)
	r.Methods("POST").Path("/string/count").Handler(countHandler)

	return r
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
