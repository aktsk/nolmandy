package receipt

import (
	"strconv"
	"time"
)

type date time.Time
type dateMS time.Time
type datePST time.Time

func (d date) MarshalJSON() ([]byte, error) {
	t := time.Time(d).Format("2006-01-02 15:04:05 Etc/GMT")
	return []byte(`"` + t + `"`), nil
}

func (d *date) UnmarshalJSON(b []byte) error {
	dateString, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02 15:04:05 Etc/GMT", dateString)
	if err != nil {
		return err
	}
	*d = date(t)

	return nil
}

func (d dateMS) MarshalJSON() ([]byte, error) {
	t := strconv.FormatInt(time.Time(d).Unix()*1000, 10)
	return []byte(`"` + t + `"`), nil
}

func (d *dateMS) UnmarshalJSON(b []byte) error {
	msString, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	sec, err := strconv.ParseInt(msString, 10, 64)
	if err != nil {
		return err
	}

	t := time.Unix(sec/1000, 0)
	*d = dateMS(t)

	return nil
}

func (d datePST) MarshalJSON() ([]byte, error) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}

	t := time.Time(d).In(loc).Format("2006-01-02 15:04:05 America/Los_Angeles")
	return []byte(`"` + t + `"`), nil
}

func (d *datePST) UnmarshalJSON(b []byte) error {
	pstString, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05 America/Los_Angeles", pstString, loc)
	if err != nil {
		return err
	}
	*d = datePST(t)

	return nil
}

// Receipt for an application
// https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
type Receipt struct {
	ReceiptType                string   `json:"receipt_type"`
	AdamID                     int64    `json:"adam_id"`
	AppItemID                  int64    `json:"app_item_id"`
	BundleID                   string   `json:"bundle_id"`
	ApplicationVersion         string   `json:"application_version"`
	DownloadID                 int64    `json:"download_id"`
	VersionExternalIdentifier  int64    `json:"version_external_identifier"`
	OriginalApplicationVersion string   `json:"original_application_version"`
	InApp                      []*InApp `json:"in_app"`
	ExpirationDate             date     `json:"-"`
	rawBundleID                []byte
	OpaqueValue                []byte `json:"-"`
	SHA1Hash                   []byte `json:"-"`
	CreationDate
	RequestDate
	OriginalPurchaseDate
}

// InApp represents the receipt for in-app purchase
// https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html#//apple_ref/doc/uid/TP40010573-CH106-SW12
type InApp struct {
	Quantity              int64  `json:"quantity,string"`
	ProductID             string `json:"product_id"`
	TransactionID         string `json:"transaction_id"`
	OriginalTransactionID string `json:"original_transaction_id"`
	WebOrderLineItemID    int64  `json:"web_order_line_item_id,omitempty"`

	IsTrialPeriod string `json:"is_trial_period"`
	ExpiresDate

	PurchaseDate
	OriginalPurchaseDate

	CancellationDate
	CancellationReason string `json:"cancellation_reason,omitempty"`
	IsInIntroPrice     bool   `json:"-"`
}

// CreationDate is the date when the app receipt was created
type CreationDate struct {
	Date    date    `json:"receipt_creation_date"`
	DateMS  dateMS  `json:"receipt_creation_date_ms"`
	DatePST datePST `json:"receipt_creation_date_pst"`
}

// RequestDate is the date when verify request was issued
type RequestDate struct {
	Date    date    `json:"request_date"`
	DateMS  dateMS  `json:"request_date_ms"`
	DatePST datePST `json:"request_date_pst"`
}

// ExpiresDate is for the subscription
type ExpiresDate struct {
	Date    date    `json:"-"`
	DateMS  dateMS  `json:"-"`
	DatePST datePST `json:"-"`
}

// PurchaseDate is the date and time that the item was purchased
type PurchaseDate struct {
	Date    date    `json:"purchase_date"`
	DateMS  dateMS  `json:"purchase_date_ms"`
	DatePST datePST `json:"purchase_date_pst"`
}

// OriginalPurchaseDate is for a transaction that restores a previous transaction, the date of the original transaction.
type OriginalPurchaseDate struct {
	Date    date    `json:"original_purchase_date"`
	DateMS  dateMS  `json:"original_purchase_date_ms"`
	DatePST datePST `json:"original_purchase_date_pst"`
}

// CancellationDate is for a transaction that was canceled by Apple customer support
type CancellationDate struct {
	Date    date    `json:"-"`
	DateMS  dateMS  `json:"-"`
	DatePST datePST `json:"-"`
}
