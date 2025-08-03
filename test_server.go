package main

import (
	"net/http"
)

func main() {

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./testfiles"))))
	http.ListenAndServe("localhost:8080", nil)
}
