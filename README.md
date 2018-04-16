# Nolmandy [![Build Status](https://travis-ci.org/aktsk/nolmandy.svg?branch=master)](https://travis-ci.org/aktsk/nolmandy)

Nolmandy is an Apple receipt processing server. You can use nolmandy instead of https://sandbox.itunes.apple.com/verifyReceipt .

Also you can use nolmandy as a receipt processing library.

**This product is at an early stage of development and not fully implemented.**

----

## Usage

### As a receipt validation command line tool

Install `nolmandy` command.

```
go get github.com/aktsk/nolmandy/cmd/nolmandy
```

Run nolmandy command to validate a receipt by Apple Root certificate.

```
cat receipt | nolmandy
```

You can validate a certificate by your own certificate.

```
cat receipt | nolmandy -certFile cert.pem
```


### As a validation server

Install `nolmandy-server` command.

```
go get github.com/aktsk/nolmandy/cmd/nolmandy-server
```

Run nolmandy server.

```
nolmandy-server -port 8000
```

Post base64 encoded receipt data to nolmandy server.

```
curl -s -H 'Content-Type:application/json' -d '{ "receipt-data": "MIIeWQYJK..." }' \
  http://localhost:8000/
```

You can use your own certificate instead of Apple certificate.

```
nolmandy-server -certFile cert.pem
```

### As a validation library

You can parse base64 encoded receipt data and validate it.

```go
package main

import (
	"log"

	"github.com/aktsk/nolmandy/receipt"
)

func main() {
	rcpt, err := receipt.ParseWithAppleRootCert("MIIT6QYJK...")
	if err != nil {
		log.Fatal("Parse error")
	}

	result, err := rcpt.Validate() // Validate() does nothing currently ...
	if err != nil {
		log.Fatal("Validation error")
	}

	if result.Status == 0 {
		log.Println("Validation success")
	}
}
```

You can use your own certificate instead of Apple root certificate like this.

```go
func main() {
	certFile, _ := os.Open("cert.pem")
	certPEM, _ := ioutil.ReadAll(certFile)
	certDER, _ := pem.Decode(certPEM)
	cert, _ = x509.ParseCertificate(certDER.Bytes)

	rcpt, err := receipt.Parse(cert, "MIIT6QYJK...")
	if err != nil {
		log.Fatal("Parse error")
	}

	result, err := rcpt.Validate() // Validate() does nothing currently ...
	if err != nil {
		log.Fatal("Validation error")
	}

	if result.Status == 0 {
		log.Println("Validation success")
	}
}
```

### Deploy nolmandy server to Google App Engine

You can run nolmandy server on Google App Engine.

```
cd appengine/app
make deploy
```

If you'd like to use your own certificate instead of Apple certificate, put a certificate file as `cert.pem` under appengine/app directory. Or you can set your certificate in app.yaml like this.

```yaml
env_variables:
  # It seems GAE/Go could not handle environment variables
  # that has return code.So use ">-" to replace return code
  # with white space
  CERTIFICATE: >-
    -----BEGIN CERTIFICATE-----
    MIIB3TCCAUagAwIBAgIEcotswjANBgkqhkiG9w0BAQsFADAoMRAwDgYDVQQKEwdB
    Y21lIENvMRQwEgYDVQQDEwtUZXN0IElzc3VlcjAgFw0xODA0MDIwNDA2MjlaGA8z
    ODQzMDQwMjA0MDYyOVowKDEQMA4GA1UEChMHQWNtZSBDbzEUMBIGA1UEAxMLVGVz
    dCBJc3N1ZXIwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMb89jkLRDjud2Xt
    tYoXscWCGTKAr+TYv7dwk7YXqabv4bhH5X15sbg9cp5UWZzu7ygoX4+N/48Aa/AZ
    Sh8ppQZYpa73ximUesF8W+ZDXeOsexoPuRyXKltjXX/eLklBldVBB8Weyip4WhmG
    TCguTUW5eFPtseIEBHxto12jB6gnAgMBAAGjEjAQMA4GA1UdDwEB/wQEAwICpDAN
    BgkqhkiG9w0BAQsFAAOBgQC+abUGkSNC5n6r4TjbCrAHZcFI0yCcK38fS2g9c7lb
    VcvltNox2SWL9oyjybdzm1iZoVtsHXuQ8RKszdVKCh7N1RUOGDgtuwfP2XnKCKoP
    W9VfLKZ+Y4YnouEZBUjsS39dgLC2EJ66e3kWfCrR6HNsSWwE0A3mVnfNUwLvgtH/
    QQ==
    -----END CERTIFICATE-----
```

----

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

----

## License

See [LICENSE](LICENSE).

----

## See Also

* [Receipt Validation Programming Guide](https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Introduction.html)
* [aktsk/kalvados: Apple receipt generator for testing](https://github.com/aktsk/kalvados)
