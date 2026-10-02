package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/ballot", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
