package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
	} else {
		certPEM := os.Getenv("CERTIFICATE")
		if certPEM != "" {
			certPEM = revertPEM(certPEM)
			certDER, _ := pem.Decode([]byte(certPEM))
			cert, err = x509.ParseCertificate(certDER.Bytes)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	http.HandleFunc("/", server.Parse(cert))
}

// It seems GAE/Go can not handle environment variables that has
// return code. So in YAML file, I  use ">-" to replace return
// code with white space. This function is for reverting back PEM data
// to original.
func revertPEM(c string) string {
	c = strings.Replace(c, " ", "\n", -1)
	c = strings.Replace(c, "\nCERTIFICATE", " CERTIFICATE", -1)
	return c
}
