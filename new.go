package adyen

import "fmt"

func NewEncryptor(url string, encryptionKey string, card Card) (*Encryptor, error) {
	e := &Encryptor{
		Url:                 url,
		EncryptionKeyString: encryptionKey,
		Card:                card,
	}

	return e, e.parseKey()
}

func NewUrl(siteKey, ver, dParam string) string {
	return fmt.Sprintf(
		"https://checkoutshopper-live-us.adyen.com/checkoutshopper/securedfields/%s/%s/securedFields.html?type=card&d=%s",
		siteKey,
		ver,
		dParam,
	)
}
