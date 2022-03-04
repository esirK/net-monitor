package main

import "net/http"

func(application *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/pings", application.handlePing)
	return mux
}
