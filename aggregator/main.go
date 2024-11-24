package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strconv"
)

func main() {

	httpListenAddr := flag.String("httpAddr", ":8080", "Listen address")
	grpcListenAddr := flag.String("grpcAddr", ":8081", "Listen address")

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	go makeGRPCTransport(*grpcListenAddr, svc)
	makeHttpTrnsport(*httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("Grpc Transport start :", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHttpTrnsport(listenAddr string, svc Aggregator) {
	fmt.Println("Http Transport start :", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.URL.Query()["obu"]
		if !ok {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing obu id"})
			return
		}
		obuID, err := strconv.Atoi(value[0])
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
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
