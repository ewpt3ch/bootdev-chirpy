package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "static"
	const filepathAssets = "assets"
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathAssets))))

	httpServe := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(httpServe.ListenAndServe())

}
