package receipt

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"
	"time"

	"github.com/guregu/null/v5"
)

func TestParseAndValidate(t *testing.T) {
	certDER, _ := pem.Decode([]byte(certificate))
	cert, err := x509.ParseCertificate(certDER.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	rcpt, err := Parse(cert, receiptData)
	if err != nil {
		t.Fatal(err)
	}

	if rcpt.ReceiptType != "ProductionSandbox" {
		t.Fatalf("Wrong receipt_type: %s", rcpt.ReceiptType)
	}

	if rcpt.BundleID != "jp.aktsk.kalvados.test" {
		t.Fatalf("Wrong bundle_id: %s", rcpt.BundleID)
	}

	creationDate := time.Unix(1518284220, 0)
	date := null.Time(rcpt.CreationDate.Date).Time
	if date.UTC() != creationDate.UTC() {
		t.Fatalf("Wrong creation_date: %v", date)
	}

	inApp := rcpt.InApp[1]

	if inApp.Quantity != 1 {
		t.Fatalf("Wrong qutantity: %d", inApp.Quantity)
	}

	if inApp.ProductID != "jp.aktsk.kalvados.test.iap1" {
		t.Fatalf("Wrong product_id: %s", inApp.ProductID)
	}

	if inApp.TransactionID != "220000359893979" {
		t.Fatalf("Wrong transaction_id: %s", inApp.TransactionID)
	}

	purchaseDate := time.Unix(1503544635, 0)
	date = null.Time(inApp.PurchaseDate.Date).Time
	if date.UTC() != purchaseDate.UTC() {
		t.Fatalf("Wrong purchase_date: %v", date)
	}

	if inApp.OriginalTransactionID != "220000348788557" {
		t.Fatalf("Wrong transaction_id: %s", inApp.OriginalTransactionID)
	}

	originalPurchaseDate := time.Unix(1500261436, 0)
	date = null.Time(inApp.OriginalPurchaseDate.Date).Time
	if date.UTC() != originalPurchaseDate.UTC() {
		t.Fatalf("Wrong original_purchase_date: %v", date)
	}

	if inApp.WebOrderLineItemID != 220000072586770 {
		t.Fatalf("Wrong web_order_line_item_id: %d", inApp.WebOrderLineItemID)
	}

	if rcpt.OriginalApplicationVersion != "49" {
		t.Fatalf("Wrong original_application_version: %s", rcpt.OriginalApplicationVersion)
	}

	originalPurchaseDate = time.Unix(1499441767, 0)
	date = null.Time(rcpt.OriginalPurchaseDate.Date).Time
	if date.UTC() != originalPurchaseDate.UTC() {
		t.Fatalf("Wrong original_purchase_date: %v", date)
	}

	if inApp.CancellationDate.Date.Valid {
		t.Fatalf("Wrong cancellation_date: %v", inApp.CancellationDate.Date)
	}

	validated, err := rcpt.Validate()
	if err != nil {
		t.Fatal(err)
	}

	if validated.Status != 0 {
		t.Fatalf("Wrong status: %d", validated.Status)
	}

	if validated.Environment != "Sandbox" {
		t.Fatalf("Wrong environment: %s", validated.Environment)
	}
}

func TestMarshalAndUnmarshalDate(t *testing.T) {
	date1 := date{}

	date1JSONString, err := json.Marshal(date1)
	if err != nil {
		t.Fatal(err)
	}

	if string(date1JSONString) != "null" {
		t.Fatalf("Wrong date1JSONString: %s", date1JSONString)
	}

	var date1JSON string
	if err := json.Unmarshal([]byte(date1JSONString), &date1JSON); err != nil {
		t.Fatal(err)
	}

	if date1JSON != "" {
		t.Fatalf("Wrong date1JSON: %s", date1JSON)
	}
}

var receiptData = `
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`

var certificate = `
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
-----END CERTIFICATE-----`
