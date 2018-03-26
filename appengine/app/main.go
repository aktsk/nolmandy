package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aktsk/nolmandy/server"
)

func init() {
	var cert *x509.Certificate

	certFile, err := os.Open("cert.pem")
	defer certFile.Close()

	if err == nil {
		certPEM, err := ioutil.ReadAll(certFile)
		if err != nil {
			log.Fatal(err)
		}

		certDER, _ := pem.Decode(certPEM)
		cert, err = x509.ParseCertificate(certDER.Bytes)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", server.Parse(cert))
}
