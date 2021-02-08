package header

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(hs Service) http.Handler {
	headerHandler := kithttp.NewServer(
		makeHeaderEndpoint(hs),
		decodeHeaderRequest,
		encodeResponse,
	)

	r := mux.NewRouter()

	r.Methods("GET", "POST").Path("/headers").Handler(headerHandler)

	return r
}

func decodeHeaderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	_, err := w.Write(response.([]byte))
	return err
}
