package server

import (
	"crypto/x509"
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
func Serve(port int, cert *x509.Certificate) {
	http.HandleFunc("/", Parse(cert))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Parse parsed receipt-data in a request
func Parse(cert *x509.Certificate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		json.NewDecoder(r.Body).Decode(&request)

		parsedReceipt := &receipt.Receipt{}

		if cert == nil {
			rcpt, err := receipt.ParseWithAppleRootCert(request.ReceiptData)
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			parsedReceipt = rcpt
		} else {
			rcpt, err := receipt.Parse(cert, request.ReceiptData)
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			parsedReceipt = rcpt
		}

		result, err := parsedReceipt.Validate()

		resultBody, err := json.Marshal(result)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(resultBody)
	}
}
