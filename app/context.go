package app

import (
	"context"

	"github.com/go-shopify/shopify"
)

type contextKey int

const (
	contextKeyShop contextKey = iota
	contextKeyLocale
	contextKeyTimestamp
)

// WithShop associates a shop to the context.
func WithShop(ctx context.Context, shop shopify.Shop) context.Context {
	return context.WithValue(ctx, contextKeyShop, shop)
}

// GetShop gets the shop associated to the specified context or nil if no shop
// is associated to the context.
func GetShop(ctx context.Context) *shopify.Shop {
	if shop, ok := ctx.Value(contextKeyShop).(shopify.Shop); ok {
		return &shop
	}

	return nil
}

// WithLocale associates a locale to the context.
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, contextKeyLocale, locale)
}

// GetLocale gets the shop associated to the specified context or nil if no shop
// is associated to the context.
func GetLocale(ctx context.Context) *string {
	if locale, ok := ctx.Value(contextKeyLocale).(string); ok {
		return &locale
	}

	return nil
}
