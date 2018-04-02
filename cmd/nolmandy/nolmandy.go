package main

import (
	"bufio"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aktsk/nolmandy/receipt"
)

func main() {
	var certFileName string

	flag.StringVar(&certFileName, "certFile", "", "Cetificate file")

	flag.Parse()

	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	receiptData := string(stdin.Bytes())

	var rcpt *receipt.Receipt
	if certFileName != "" {
		certFile, err := os.Open(certFileName)
		if err != nil {
			handleError(err)
		}

		defer certFile.Close()

		certPEM, err := ioutil.ReadAll(certFile)
		if err != nil {
			handleError(err)
		}

		certDER, _ := pem.Decode(certPEM)
		cert, err := x509.ParseCertificate(certDER.Bytes)
		if err != nil {
			handleError(err)
		}

		rcpt, err = receipt.Parse(cert, receiptData)
		if err != nil {
			handleError(err)
		}
	} else {
		var err error
		rcpt, err = receipt.ParseWithAppleRootCert(receiptData)
		if err != nil {
			handleError(err)
		}
	}

	res, err := rcpt.Validate()

	json, err := json.Marshal(res)
	if err != nil {
		handleError(err)
	}

	fmt.Println(string(json))
}

func handleError(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
