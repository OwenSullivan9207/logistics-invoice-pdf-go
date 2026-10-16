package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	c := infraiClient{key: key, http: &http.Client{}}
	http.HandleFunc("POST /invoices", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var o Order
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&o); err != nil {
			http.Error(w, "invalid order JSON", http.StatusBadRequest)
			return
		}
		result, err := sendInvoice(r.Context(), c, o)
		if err != nil {
			http.Error(w, err.Error(), clientStatus(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
