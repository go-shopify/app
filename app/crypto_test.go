package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-shopify/shopify"
)

func TestNewHMACHandler(t *testing.T) {
	apiSecret := shopify.APISecret("abcdefgh")

	handler := NewHMACHandler(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
		apiSecret,
	)

	t.Run("good", func(t *testing.T) {
		w := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}

		req, _ := http.NewRequest(http.MethodGet, "http://myhost?code=0907a61c0c8d55e99db179b68161bc00&hmac=26015c6ad20dccdc7017bc4ad3b7c7b239a18db8c79e26945e13d0ca551ae996&timestamp=1337178173&shop=some-shop.myshopify.com&state=0.6784241404160823", nil)

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, w.Code)
			t.Errorf("Body follows:\n%s", w.Body.String())
		}
	})

	t.Run("bad", func(t *testing.T) {
		w := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}

		req, _ := http.NewRequest(http.MethodGet, "http://myhost?code=0907a61c0c8d55e99db179b68161bc00&hmac=26015c6ad20dccdc7017bc4ad3b7c7b239a18db8c79e26945e13d0ca551ae996&shop=some-shop.myshopify.com&state=0.6784241404160823&timestamp=2337178173", nil)

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected %d but got %d", http.StatusForbidden, w.Code)
			t.Errorf("Body follows:\n%s", w.Body.String())
		}
	})

	t.Run("missing", func(t *testing.T) {
		w := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}

		req, _ := http.NewRequest(http.MethodGet, "http://myhost?code=0907a61c0c8d55e99db179b68161bc00&shop=some-shop.myshopify.com&state=0.6784241404160823&timestamp=2337178173", nil)

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected %d but got %d", http.StatusBadRequest, w.Code)
			t.Errorf("Body follows:\n%s", w.Body.String())
		}
	})

	t.Run("encoding-encoded", func(t *testing.T) {
		w := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}

		req, _ := http.NewRequest(http.MethodGet, "http://myhost?code=0907a61c0c8d55e99db179b68161bc00&hmac=02efe8e8d299cb8d403022a485b6685cd186f499258110568c9af306a3598083&timestamp=1337178173&shop=some-shop.myshopify.com&state=a%3D%3D", nil)

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, w.Code)
			t.Errorf("Body follows:\n%s", w.Body.String())
		}
	})

	t.Run("encoding-decoded", func(t *testing.T) {
		w := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}

		req, _ := http.NewRequest(http.MethodGet, "http://myhost?code=0907a61c0c8d55e99db179b68161bc00&hmac=02efe8e8d299cb8d403022a485b6685cd186f499258110568c9af306a3598083&timestamp=1337178173&shop=some-shop.myshopify.com&state=a==", nil)

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, w.Code)
			t.Errorf("Body follows:\n%s", w.Body.String())
		}
	})
}
