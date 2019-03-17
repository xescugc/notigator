package service

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(s Service) http.Handler {
	getSourcesHandler := kithttp.NewServer(
		makeGetSources(s),
		decodeGetSourcesRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)

	getSourceNotificationsHandler := kithttp.NewServer(
		makeGetSourceNotifications(s),
		decodeGetSourceNotificationsRequest,
		encodeResponse,
		kithttp.ServerErrorEncoder(encodeError),
	)

	r := mux.NewRouter()

	r.Handle("/api/sources", getSourcesHandler).Methods("GET")
	r.Handle("/api/sources/{sourceID}/notifications", getSourceNotificationsHandler).Methods("GET")

	return r
}

func decodeGetSourcesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetSourceNotificationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getSourceNotificationsRequest{
		SourceID: mux.Vars(r)["sourceID"],
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
