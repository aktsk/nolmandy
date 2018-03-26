# Nolmandy

Nolmandy is an Apple receipt processing server. You can use nolmandy instead of https://sandbox.itunes.apple.com/verifyReceipt .

Also you can use nolmandy as a receipt processing library.

**This product is at an early stage of development and not fully implemented.**

----

## Usage

### As a validation server

Compile nolmandy.

```
make
```

Run nolmandy server.

```
bin/nolmandy-server -port 8000
```

Post base64 encoded receipt data to nolmandy server.

```
curl -s -H 'Content-Type:application/json' -d '{ "receipt-data": "MIIeWQYJK..." }' \
  http://localhost:8000/
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

### Deploy nolmandy server to Google App Engine

You can run nolmandy server on Google App Engine.

```
cd appengine/app
make deploy
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

