package token

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Token struct {
	keyset jwk.Set
}

func NewToken(keyset jwk.Set) *Token {
	return &Token{keyset: keyset}
}

func (t *Token) Verify(tokenStr string) (jwt.Token, error) {
	parsedToken, err := jwt.ParseString(
		tokenStr,
		jwt.WithKeySet(t.keyset),
		jwt.WithValidate(true),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return parsedToken, nil
}
