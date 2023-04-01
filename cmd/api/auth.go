package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type auth struct {
	Issuer        string `json:"issuer"`
	Audience      string `json:"audience"`
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type tokenPairs struct {
	Token       string `json:"access_token"` // The actual JWT token
	RefresToken string `json:"refresh_token"`
}

type claims struct {
	jwt.RegisteredClaims
}

func (a *auth) GenerateTokenPair(user *jwtUser) (tokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = a.Audience
	claims["iss"] = a.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(a.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		log.Println("signing access token: ", err)
		return tokenPairs{}, err
	}

	// Create a refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set the expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(a.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(a.Secret))
	if err != nil {
		log.Println("signing refresh token: ", err)
		return tokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	pair := tokenPairs{
		Token:       signedAccessToken,
		RefresToken: signedRefreshToken,
	}

	// Return TokenPairs
	return pair, nil
}

func (a *auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     a.CookieName,
		Path:     a.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(a.RefreshExpiry),
		MaxAge:   int(a.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   a.CookieDomain,
		Secure:   true,
		HttpOnly: true,
	}
}

func (a *auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     a.CookieName,
		Path:     a.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   a.CookieDomain,
		Secure:   true,
		HttpOnly: true,
	}
}

func (a *auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *claims, error) {
	// To avoid caching of requests that include this header
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	// Sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// Check for "Bearer XXX
	headerParts := strings.Fields(authHeader)
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// Check for the word bearer
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]
	claims := &claims{}

	// Parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	if claims.Issuer != a.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}
