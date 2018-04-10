// +build ignore

package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/aktsk/nolmandy/server"
)

func main() {
	var port int
	var certFileName string

	flag.IntVar(&port, "port", 8000, "Port to listen")
	flag.StringVar(&certFileName, "certFile", "", "Certificate file")

	flag.Parse()

	var cert *x509.Certificate

	if certFileName != "" {
		certFile, err := os.Open(certFileName)
		if err != nil {
			log.Fatal(err)
		}

		defer certFile.Close()

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

	server.Serve(port, cert)
}
