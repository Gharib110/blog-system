package api

import (
	"fmt"
	"net/http"
)

// enableCORS use for enabling the cors for our client
func (rcf *RestConf) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}
