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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
