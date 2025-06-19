package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	server := http.Server{
		Handler: mux,
		Addr: ":9000",
	}

	server.ListenAndServe()
}
