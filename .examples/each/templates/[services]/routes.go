package {{ .Each.Item }}

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/{{ .Each.Item }}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{{ .Each.Item }} service"))
	})
}
