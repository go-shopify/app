package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-shopify/shopify"
)

func TestAPIHandler(t *testing.T) {
	stok := &sessionToken{
		Shop: "myshop",
		OAuthToken: shopify.OAuthToken{
			AccessToken: "abc",
			Scope: shopify.Scope{
				shopify.PermissionReadProducts,
			},
		},
	}

	ctx := context.Background()
	oauthTokenStorage := &MemoryOAuthTokenStorage{}
	oauthTokenStorage.UpdateOAuthToken(ctx, stok.Shop, stok.OAuthToken)

	handler := NewAPIHandler(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			shop, ok := shopify.GetShop(req.Context())

			if !ok {
				t.Fatalf("expected true")
			}

			if shop != stok.Shop {
				t.Errorf("expected `%s` but got `%s`", shop, stok.Shop)
			}

			w.WriteHeader(http.StatusOK)
		}),
		oauthTokenStorage,
	)

	t.Run("no cookie", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("valid cookie missing shop", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		req.AddCookie(stok.AsCookie())
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("valid cookie", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		req.AddCookie(&http.Cookie{Name: shopifyShopCookieName, Value: string(stok.Shop)})
		req.AddCookie(stok.AsCookie())
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("invalid cookie (base64)", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		cookie := &http.Cookie{
			Name:  sessionTokenCookieName,
			Value: "$",
		}
		req.AddCookie(cookie)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("invalid cookie (json)", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		cookie := &http.Cookie{
			Name:  sessionTokenCookieName,
			Value: "blah",
		}
		req.AddCookie(cookie)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("no oauth token available", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		stok2 := &sessionToken{
			Shop: "myshop2",
			OAuthToken: shopify.OAuthToken{
				AccessToken: "abc",
				Scope: shopify.Scope{
					shopify.PermissionReadProducts,
				},
			},
		}
		req.AddCookie(stok2.AsCookie())
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("altered session token", func(t *testing.T) {
		w := &httptest.ResponseRecorder{}
		req := httptest.NewRequest(http.MethodGet, "https://foo", nil)
		stok2 := &sessionToken{
			Shop: "myshop",
			OAuthToken: shopify.OAuthToken{
				AccessToken: "abd",
				Scope: shopify.Scope{
					shopify.PermissionReadProducts,
				},
			},
		}
		req.AddCookie(stok2.AsCookie())
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d but got %d", http.StatusForbidden, w.Code)
		}
	})
}
