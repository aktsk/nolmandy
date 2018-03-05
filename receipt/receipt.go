package receipt

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"io/ioutil"
	"time"

	_ "github.com/aktsk/nolmandy/statik" // Need to load assets
	"github.com/fullsailor/pkcs7"
	"github.com/rakyll/statik/fs"
)

// Result is the validation result
type Result struct {
	Status            int      `json:"status"`
	Environment       string   `json:"environment"`
	Receipt           *Receipt `json:"receipt"`
	LatestReceiptInfo []InApp  `json:"latest_receipt_info,omitempty"`
	LatestReceipt     string   `json:"latest_receipt,omitempty"`
	IsRetryable       bool     `json:"is-retryable,omitempty"`
}

// ParseWithAppleRootCert parses base 64 encoded receipt data with Apple Inc Root Certificate
func ParseWithAppleRootCert(data string) (*Receipt, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	rootCertFile, err := statikFS.Open("/AppleIncRootCertificate.cer")
	if err != nil {
		return nil, err
	}

	rootCertBytes, err := ioutil.ReadAll(rootCertFile)
	if err != nil {
		return nil, err
	}

	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return nil, err
	}

	return Parse(rootCert, data)
}

// Parse parsed base 64 encoded receipt data with a given certificate
func Parse(root *x509.Certificate, data string) (*Receipt, error) {
	receiptData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	pkcs, err := pkcs7.Parse(receiptData)
	if err != nil {
		return nil, err
	}

	if err := verifySignerCert(root, pkcs); err != nil {
		return nil, err
	}

	if err := pkcs.Verify(); err != nil {
		return nil, err
	}

	receipt, err := parsePKCS(pkcs)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// Validate is for validating receipt
func (r *Receipt) Validate() (Result, error) {
	return Result{
		Status:      0,
		Environment: "Sandbox",
		Receipt:     r,
	}, nil
}

func verifySignerCert(root *x509.Certificate, pkcs *pkcs7.PKCS7) error {
	roots := x509.NewCertPool()
	roots.AddCert(root)

	signer := pkcs.GetOnlySigner()

	intermediates := x509.NewCertPool()
	for _, cert := range pkcs.Certificates {
		if cert != signer && !cert.Equal(root) {
			intermediates.AddCert(cert)
		}
	}

	_, err := signer.Verify(x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:         roots,
	})

	if err != nil {
		return err
	}

	return nil
}

type attribute struct {
	Type    int
	Version int
	Value   []byte
}

func parsePKCS(pkcs *pkcs7.PKCS7) (*Receipt, error) {
	var receipt Receipt

	var r asn1.RawValue
	_, err := asn1.Unmarshal(pkcs.Content, &r)
	if err != nil {
		return nil, err
	}
	rest := r.Bytes
	for len(rest) > 0 {
		var ra attribute
		rest, err = asn1.Unmarshal(rest, &ra)
		if err != nil {
			return nil, err
		}
		switch ra.Type {
		case 2:
			if _, err = asn1.Unmarshal(ra.Value, &receipt.BundleID); err != nil {

				return nil, err
			}
			receipt.rawBundleID = ra.Value
		case 3:
			if _, err = asn1.Unmarshal(ra.Value, &receipt.ApplicationVersion); err != nil {
				return nil, err
			}
		case 4:
			receipt.OpaqueValue = ra.Value
		case 5:
			receipt.SHA1Hash = ra.Value
		case 12:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			receipt.CreationDate.Date = date(t)
			receipt.CreationDate.DateMS = dateMS(t)
			receipt.CreationDate.DatePST = datePST(t)

		case 17:
			var inApp *InApp
			inApp, err = parseInApp(ra.Value)
			if err != nil {
				return nil, err
			}
			receipt.InApp = append(receipt.InApp, inApp)
		case 19:
			if _, err = asn1.Unmarshal(ra.Value, &receipt.OriginalApplicationVersion); err != nil {
				return nil, err
			}
		case 21:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			receipt.ExpirationDate = date(t)

		// Field types below are not listed in https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
		case 0:
			if _, err = asn1.Unmarshal(ra.Value, &receipt.ReceiptType); err != nil {
				return nil, err
			}
		case 18:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			receipt.OriginalPurchaseDate.Date = date(t)
			receipt.OriginalPurchaseDate.DateMS = dateMS(t)
			receipt.OriginalPurchaseDate.DatePST = datePST(t)
		}
	}

	loc, _ := time.LoadLocation("Etc/GMT")
	now := time.Now().In(loc)
	receipt.RequestDate.Date = date(now)
	receipt.RequestDate.DateMS = dateMS(now)
	receipt.RequestDate.DatePST = datePST(now)

	return &receipt, nil
}

func parseInApp(data []byte) (*InApp, error) {
	var inApp InApp

	var r asn1.RawValue
	_, err := asn1.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	data = r.Bytes
	for len(data) > 0 {
		var ra attribute
		data, err = asn1.Unmarshal(data, &ra)
		if err != nil {
			return nil, err
		}
		switch ra.Type {
		case 1701:
			if _, err = asn1.Unmarshal(ra.Value, &inApp.Quantity); err != nil {
				return nil, err
			}
		case 1702:
			if _, err = asn1.Unmarshal(ra.Value, &inApp.ProductID); err != nil {
				return nil, err
			}
		case 1703:
			if _, err = asn1.Unmarshal(ra.Value, &inApp.TransactionID); err != nil {
				return nil, err
			}
		case 1704:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			inApp.PurchaseDate.Date = date(t)
			inApp.PurchaseDate.DateMS = dateMS(t)
			inApp.PurchaseDate.DatePST = datePST(t)
		case 1705:
			if _, err = asn1.Unmarshal(ra.Value, &inApp.OriginalTransactionID); err != nil {
				return nil, err
			}
		case 1706:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			inApp.OriginalPurchaseDate.Date = date(t)
			inApp.OriginalPurchaseDate.DateMS = dateMS(t)
			inApp.OriginalPurchaseDate.DatePST = datePST(t)
		case 1708:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			inApp.ExpiresDate.Date = date(t)
		case 1711:
			if _, err = asn1.Unmarshal(ra.Value, &inApp.WebOrderLineItemID); err != nil {
				return nil, err
			}
		case 1712:
			t, err := asn1ParseTime(ra.Value)
			if err != nil {
				return nil, err
			}
			inApp.CancellationDate.Date = date(t)
		case 1719:
			var intRoprice int
			if _, err = asn1.Unmarshal(ra.Value, &intRoprice); err != nil {
				return nil, err
			}
			inApp.IsInIntroPrice = (intRoprice != 0)
		}
	}

	if inApp.IsTrialPeriod == "" {
		inApp.IsTrialPeriod = "false"
	}

	return &inApp, nil
}

func asn1ParseTime(data []byte) (time.Time, error) {
	var str string
	if _, err := asn1.Unmarshal(data, &str); err != nil {
		return time.Time{}, err
	}
	if str == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, str)
}
