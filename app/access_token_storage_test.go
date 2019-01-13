package app

import (
	"context"
	"testing"

	"github.com/go-shopify/shopify"
)

func TestMemoryAccessTokenStorage(t *testing.T) {
	var storage AccessTokenStorage = &MemoryAccessTokenStorage{}

	ctx := context.Background()
	shop := shopify.Shop("myshop")

	accessToken, err := storage.GetAccessToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if accessToken != "" {
		t.Errorf("expected no access token: %s", accessToken)
	}

	ref := shopify.AccessToken("token")
	err = storage.UpdateAccessToken(ctx, shop, ref)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	accessToken, err = storage.GetAccessToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if accessToken != ref {
		t.Errorf("expected a different access token: %s", accessToken)
	}

	err = storage.DeleteAccessToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	accessToken, err = storage.GetAccessToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if accessToken != "" {
		t.Errorf("expected no access token: %s", accessToken)
	}
}
