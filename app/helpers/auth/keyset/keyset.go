package keyset

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type KeySet struct {
	jwkURL            string
	minRefreshMinutes int64
}

func NewKeySetWithCache(jwkURL string, minRefreshMinutes int64) *KeySet {
	return &KeySet{jwkURL: jwkURL, minRefreshMinutes: minRefreshMinutes}
}

func (k KeySet) NewKeySet() (jwk.Set, error) {
	jwkCache := jwk.NewCache(context.Background())

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
	err := jwkCache.Register(k.jwkURL, jwk.WithMinRefreshInterval(time.Duration(k.minRefreshMinutes)*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("failed to register JWK URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, k.jwkURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK: %w", err)
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, k.jwkURL), nil
}
