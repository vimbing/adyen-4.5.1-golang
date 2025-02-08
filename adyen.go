package adyen

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	FIELD_NUMBER       = "number"
	FIELD_EXPIRY_MONTH = "expMonth"
	FIELD_EXPIRY_YEAR  = "expYear"
	FIELD_CVV          = "cvv"
)

// func (en *Encryptor) to(e string) string {
// 	key := base64.RawStdEncoding.EncodeToString([]byte(e))

// 	re1 := regexp.MustCompile(`=`)
// 	re2 := regexp.MustCompile(`\+`)
// 	re3 := regexp.MustCompile(`\/`)

// 	key = re1.ReplaceAllString(key, "")
// 	key = re2.ReplaceAllString(key, "-")
// 	key = re3.ReplaceAllString(key, "_")

// 	return key
// }

func (en *Encryptor) ro(e string) string {
	if len(e)%2 == 1 {
		e = "0" + e
	}

	t, _ := hex.DecodeString(e)
	return string(t)
}

func (e *Encryptor) parseKey() error {
	r := strings.Split(e.EncryptionKeyString, "|")

	n := r[0]
	o := r[1]

	i := e.ro(n)
	a := e.ro(o)

	eBytes := []byte(i)
	eInt := new(big.Int).SetBytes(eBytes)

	nBytes := []byte(a)
	nInt := new(big.Int).SetBytes(nBytes)

	publicKey := &rsa.PublicKey{
		N: nInt,
		E: int(eInt.Int64()),
	}

	key, err := jwk.FromRaw(publicKey)

	if err != nil {
		return err
	}

	if err := key.Set(jwk.KeyIDKey, "asf-key"); err != nil {
		return err
	}

	e.EncryptionKey = key

	return nil
}

func (e *Encryptor) encrypt(fieldName, value string) (string, error) {
	data := map[string]string{}

	switch fieldName {
	case FIELD_NUMBER:
		data = map[string]string{
			"number":                value,
			"activate":              "3",
			"deactivate":            "1",
			"generationtime":        e.genTime,
			"numberBind":            "1",
			"numberFieldBlurCount":  "1",
			"numberFieldClickCount": "1",
			"numberFieldFocusCount": "3",
			"numberFieldKeyCount":   "2",
			"numberFieldLog":        "fo@5956,cl@5960,bl@5973,fo@6155,fo@6155,Md@6171,KL@6173,pa@6173",
			"numberFieldPasteCount": "1",
			"referrer":              e.Url,
		}
	case FIELD_EXPIRY_MONTH:
		data = map[string]string{
			"expiryMonth":    value,
			"generationtime": e.genTime,
		}
	case FIELD_EXPIRY_YEAR:
		data = map[string]string{
			"expiryYear":     value,
			"generationtime": e.genTime,
		}
	case FIELD_CVV:
		data = map[string]string{
			"activate":           "1",
			"cvc":                value,
			"cvcBind":            "1",
			"cvcFieldClickCount": "1",
			"cvcFieldFocusCount": "2",
			"cvcFieldKeyCount":   "4",
			"cvcFieldLog":        "fo@20328,fo@20328,cl@20329,KN@20344,KN@20347,KN@20349,KN@20351",
			"generationtime":     e.genTime,
			"referrer":           e.Url,
		}
	}

	marshaledData, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	headers := jwe.NewHeaders()

	headers.Set("alg", "RSA-OAEP")
	headers.Set("enc", "A256CBC-HS512")
	headers.Set("version", "1")

	encrypted, err := jwe.Encrypt(
		marshaledData,
		jwe.WithKey(jwa.RSA_OAEP, e.EncryptionKey),
		jwe.WithProtectedHeaders(headers),
		jwe.WithContentEncryption(jwa.A256CBC_HS512),
	)

	if err != nil {
		return "", err
	}

	return string(encrypted), nil

}

func (e *Encryptor) Encrypt() (*Card, error) {
	e.genTime = time.Now().Format("2006-01-02T15:04:05Z")

	cardNumber, err := e.encrypt(FIELD_NUMBER, e.Card.Number)

	if err != nil {
		return nil, err
	}

	expiryMonth, err := e.encrypt(FIELD_EXPIRY_MONTH, e.Card.ExpiryMonth)

	if err != nil {
		return nil, err
	}

	expiryYear, err := e.encrypt(FIELD_EXPIRY_YEAR, e.Card.ExpiryYear)

	if err != nil {
		return nil, err
	}

	cvv, err := e.encrypt(FIELD_CVV, e.Card.CVV)

	if err != nil {
		return nil, err
	}

	return &Card{
		ExpiryMonth: expiryMonth,
		Number:      cardNumber,
		ExpiryYear:  expiryYear,
		CVV:         cvv,
	}, nil
}
