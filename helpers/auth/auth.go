package auth

import (
	"goBoard/internal/core/domain"
	"log"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	accessTokenCookieName = "access-token"
	jwtSecretKey          = "some-secret-key"
)

func GenerateTokenAndSetCookies(member *domain.Member, w http.ResponseWriter, timeout time.Duration) error {
	accessToken, exp, err := generateAccessToken(member, timeout)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, w)
	setUserCookie(member, exp, w)

	return nil
}

func generateAccessToken(member *domain.Member, timeout time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(1 * timeout)

	return generateToken(member, exp, []byte(jwtSecretKey))
}

func generateToken(member *domain.Member, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	t := jwt.New()
	t.Set(jwt.SubjectKey, `goBoard`)
	t.Set(jwt.AudienceKey, `goBoard Users`)
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, expirationTime)
	t.Set("member", member.Name)
	//
	//key, err := rsa.GenerateKey(bytes.NewReader(secret), 2048)
	//if err != nil {
	//	log.Printf("failed to generate private key: %s", err)
	//	return "", time.Now(), err
	//}

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.HS256, secret))
	if err != nil {
		log.Printf("failed to sign token: %s", err)
		return "", time.Now(), err
	}

	return string(signed), expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, w http.ResponseWriter) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	http.SetCookie(w, cookie)
}

func setUserCookie(member *domain.Member, expiration time.Time, w http.ResponseWriter) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = member.Name
	cookie.Expires = expiration
	cookie.Path = "/"

	http.SetCookie(w, cookie)
}
