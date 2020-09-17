package main

import (
	"encoding/json"
	"net/http"
)

// TestResp ...
type TestResp struct {
	Msg       string `json:"message"`
	URL       string `json:"url"`
	UserAgent string `json:"user-agent"`
	Method    string `json:"http-method"`
	IP        string `json:"IP"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := &TestResp{}
		resp.Msg = "Hello world SaniyaShenpion"
		resp.URL = r.URL.String()
		resp.UserAgent = r.UserAgent()
		resp.Method = r.Method
		resp.IP = r.RemoteAddr

		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "JSON marshal error", http.StatusInternalServerError)
		}

		w.Write(b)
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
