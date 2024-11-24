package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"net/http"
)

func main() {

	listenAddr := flag.String("listen", ":8080", "Listen address")

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	makeHttpTrnsport(*listenAddr, svc)
}

func makeHttpTrnsport(listenAddr string, svc Aggregator) {
	fmt.Println("Http Transport start :", listenAddr)
	http.HandleFunc("/aggregate", HandleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func HandleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}
	}
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
