package jwt

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Parser struct
type Parser struct {
	Config Config
}

// Parse tries to extract and validate token from request.
// See "Config.TokenLookup" for possible ways to pass token in request.
func (jp *Parser) Parse(r *http.Request) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(jp.Config.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = jp.jwtFromHeader(r, parts[1])
	case "query":
		token, err = jp.jwtFromQuery(r, parts[1])
	case "cookie":
		token, err = jp.jwtFromCookie(r, parts[1])
	}

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(jp.Config.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return jp.Config.Secret, nil
	})
}

func (jp *Parser) jwtFromHeader(r *http.Request, key string) (string, error) {
	authHeader := r.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (jp *Parser) jwtFromQuery(r *http.Request, key string) (string, error) {
	token := r.URL.Query().Get(key)

	if token == "" {
		return "", errors.New("Query token empty")
	}

	return token, nil
}

func (jp *Parser) jwtFromCookie(r *http.Request, key string) (string, error) {
	cookie, _ := r.Cookie(key)

	if nil == cookie {
		return "", errors.New("Cookie token empty")
	}

	return cookie.Value, nil
}
