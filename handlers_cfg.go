package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
)

func (cfg *apiConfig) mwMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) reqCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, strconv.Itoa(int(cfg.fileserverHits.Load())))

	w.Write([]byte(body))
}

func (cfg *apiConfig) resetCount(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hit count reset to 0"))
}
