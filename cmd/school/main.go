package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "В голове было смешнее")
		fmt.Fprintln(w, "Request url: "+r.URL.String())
		fmt.Fprintln(w, "User agent: "+r.UserAgent())
		fmt.Fprintln(w, "Method: "+r.Method)
		fmt.Fprintln(w, "IP: "+r.RemoteAddr)
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
