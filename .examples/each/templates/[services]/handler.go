package {{ .Each.Item }}

import "net/http"

// Handler returns the HTTP handler for the {{ .Each.Item }} service.
func Handler() http.Handler {
	mux := http.NewServeMux()
	registerRoutes(mux)
	return mux
}
