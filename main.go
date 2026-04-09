package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "static"
	const filepathApp = "/app/"
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle(filepathApp, http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

	httpServe := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(httpServe.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
