package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aktsk/nolmandy/receipt"
)

// Request is for request to a receipt validation server
type Request struct {
	ReceiptData string `json:"receipt-data"`
	Password    string `json:"password"`
}

// Serve is for serving receipt vefirification
func Serve(port int) {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var request Request
	json.NewDecoder(r.Body).Decode(&request)

	rcpt, err := receipt.ParseWithAppleRootCert(request.ReceiptData)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := rcpt.Validate()

	resultBody, err := json.Marshal(result)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resultBody)
}
