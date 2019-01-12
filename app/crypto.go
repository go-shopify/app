package app

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-shopify/shopify"
)

const stateSize = 16

func generateRandomState() (string, error) {
	var stateData [stateSize]byte

	if _, err := rand.Read(stateData[:]); err != nil {
		return "", fmt.Errorf("could not generate random state: %s", err)
	}

	return base64.URLEncoding.EncodeToString(stateData[:]), nil
}

func computeHMAC(values url.Values, apiSecret shopify.APISecret) string {
	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(values.Encode()))
	return hex.EncodeToString(hmac.Sum(nil))
}

func verifyHMAC(hmac string, values url.Values, apiSecret shopify.APISecret) error {
	expected := computeHMAC(values, apiSecret)

	if hmac != expected {
		return fmt.Errorf("HMAC verification failed: expected `%s` but got `%s`", expected, hmac)
	}

	return nil
}

// NewHMACHandler wraps an existing handler and adds HMAC verification logic.
func NewHMACHandler(handler http.Handler, apiSecret shopify.APISecret) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()

		hmac := values.Get("hmac")

		if hmac == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing `hmac` parameter.")
			return
		}

		values.Del("hmac")

		if err := verifyHMAC(hmac, values, apiSecret); err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "HMAC verification failed.")
			return
		}

		handler.ServeHTTP(w, req)
	})
}
