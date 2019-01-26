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
	s := values.Encode()
	s, _ = url.QueryUnescape(s)

	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(s))
	return hex.EncodeToString(hmac.Sum(nil))
}

func verifyHMAC(h string, values url.Values, apiSecret shopify.APISecret) error {
	expected := computeHMAC(values, apiSecret)

	if !hmac.Equal([]byte(h), []byte(expected)) {
		return fmt.Errorf("HMAC verification failed: expected `%s` but got `%s`", expected, h)
	}

	return nil
}

func injectHMAC(values url.Values, apiSecret shopify.APISecret) {
	hmac := computeHMAC(values, apiSecret)
	values.Set("hmac", hmac)
}

func injectSignature(values url.Values, apiSecret shopify.APISecret) {
	signature := computeHMAC(values, apiSecret)
	values.Set("signature", signature)
}

// newHMACHandler wraps an existing handler and adds HMAC verification logic.
func newHMACHandler(handler http.Handler, apiSecret shopify.APISecret) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()

		hmac := values.Get("hmac")

		if hmac == "" {
			if hmac = values.Get("signature"); hmac == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing `hmac` or `signature` parameter.")
				return
			}
		}

		values.Del("hmac")
		values.Del("signature")

		if err := verifyHMAC(hmac, values, apiSecret); err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "HMAC verification failed.")
			return
		}

		handler.ServeHTTP(w, req)
	})
}
