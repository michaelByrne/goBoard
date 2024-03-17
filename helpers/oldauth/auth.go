package auth

import (
	v5claims "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"goBoard/internal/core/domain"
	"net/http"
	"time"
)

const (
	accessTokenCookieName = "access-token"
	jwtSecretKey          = "some-secret-key"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

type Claims struct {
	Name string `json:"name"`
	v5claims.RegisteredClaims
}

func GetJWTClaims(c echo.Context) v5claims.Claims {
	return Claims{}
}

func GenerateTokensAndSetCookies(member *domain.Member, c echo.Context, timeout time.Duration) error {
	accessToken, exp, err := generateAccessToken(member, timeout)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(member, exp, c)

	return nil
}

func generateAccessToken(member *domain.Member, timeout time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * timeout)

	return generateToken(member, expirationTime, []byte(GetJWTSecret()))
}

func generateToken(member *domain.Member, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		Name: member.Name,
		RegisteredClaims: v5claims.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: v5claims.NewNumericDate(expirationTime),
		},
	}

	token := v5claims.NewWithClaims(v5claims.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(member *domain.Member, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = member.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func JWTErrorChecker(c echo.Context, err error) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("login"))
}
