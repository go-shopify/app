package shopify

import (
	"context"
	"testing"
)

func TestWithShop(t *testing.T) {
	ctx := context.Background()
	shop := Shop("myshop")

	_, ok := GetShop(ctx)

	if ok {
		t.Errorf("expected false")
	}

	ctx = WithShop(ctx, shop)

	value, ok := GetShop(ctx)

	if !ok {
		t.Errorf("expected true")
	}

	if value != shop {
		t.Errorf("expected `%s` but got `%s`", shop, value)
	}
}

func TestWithAccessToken(t *testing.T) {
	ctx := context.Background()
	accessToken := AccessToken("abc")

	_, ok := GetAccessToken(ctx)

	if ok {
		t.Errorf("expected false")
	}

	ctx = WithAccessToken(ctx, accessToken)

	value, ok := GetAccessToken(ctx)

	if !ok {
		t.Errorf("expected true")
	}

	if value != accessToken {
		t.Errorf("expected `%s` but got `%s`", accessToken, value)
	}
}
