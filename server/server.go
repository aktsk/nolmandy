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

// Serve is for serving receipt verification
func Serve(port int, cert *x509.Certificate) {
	http.HandleFunc("/", Parse(cert))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Parse parsed receipt-data in a request
func Parse(cert *x509.Certificate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		json.NewDecoder(r.Body).Decode(&request)

		var rcpt *receipt.Receipt
		var result receipt.Result
		var err error

		if cert == nil {
			cert, err = receipt.GetAppleRootCert()
			if err != nil {
				log.Print(err)
				result.Status = 21100
				result.IsRetryable = true
				writeResult(w, result)
				return
			}
		}

		rcpt, err = receipt.Parse(cert, request.ReceiptData)

		if err != nil {
			log.Print(err)
			result.Status = 21002
		} else {
			result, err = rcpt.Validate()
		}

		writeResult(w, result)
	}
}

func writeResult(w http.ResponseWriter, result receipt.Result) {
	resultBody, err := json.Marshal(result)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resultBody)
}
