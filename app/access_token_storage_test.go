package app

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-shopify/shopify"
)

func TestMemoryOAuthTokenStorage(t *testing.T) {
	var storage OAuthTokenStorage = &MemoryOAuthTokenStorage{}

	ctx := context.Background()
	shop := shopify.Shop("myshop")

	oauthToken, err := storage.GetOAuthToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if oauthToken != nil {
		t.Errorf("expected no OAuth token: %s", oauthToken)
	}

	ref := shopify.OAuthToken{
		AccessToken: "token",
		Scope:       shopify.Scope{shopify.PermissionWriteProducts},
	}

	err = storage.UpdateOAuthToken(ctx, shop, ref)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	oauthToken, err = storage.GetOAuthToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if !reflect.DeepEqual(*oauthToken, ref) {
		t.Errorf("expected a different OAuth token: %s", *oauthToken)
	}

	err = storage.DeleteOAuthToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	oauthToken, err = storage.GetOAuthToken(ctx, shop)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if oauthToken != nil {
		t.Errorf("expected no OAuth token: %s", oauthToken)
	}
}
