package main

import (
	"fmt"
	"net/http"
	{{ range .Scaffold.services }}
	"app/{{ . }}"{{ end }}
)

func main() {
	mux := http.NewServeMux()
	{{ range .Scaffold.services }}
	mux.Handle("/{{ . }}/", {{ . }}.Handler()){{ end }}

	fmt.Println("Starting server with {{ .Computed.service_count }} services...")
	http.ListenAndServe(":8080", mux)
}
