package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hzxiao/goutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrEmptyToken = errors.New("token: empty token")
)

type Context interface {
	ToMap() goutil.Map
	LoadFromMap(data map[string]interface{}) error
}

var conf *Config

type Config struct {
	Secret       string
	LookupMethod []string
}

func Init(c *Config) {
	conf = c
}

func ParseRequest(r *http.Request, v Context) (err error) {
	token, err := getToken(r)
	if err != nil {
		return
	}
	return Parse(token, conf.Secret, v)
}

func getToken(r *http.Request) (string, error) {
	var token string
	var err error
	for _, method := range conf.LookupMethod {
		if token != "" {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), "-")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = jwtFromHeader(r, v)
		case "query":
			token, err = jwtFromQuery(r, v)
		case "cookie":
			token, err = jwtFromCookie(r, v)
		}
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func jwtFromHeader(r *http.Request, key string) (string, error) {
	auth := r.Header.Get(key)
	if auth == "" {
		return "", ErrEmptyToken
	}

	var t string
	// Parse the header to get the token part.
	fmt.Sscanf(auth, "Bearer %s", &t)
	if t == "" {
		return "", ErrEmptyToken
	}
	return t, nil
}

func jwtFromQuery(r *http.Request, key string) (string, error) {
	if values, ok := r.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], nil
	}

	return "", ErrEmptyToken
}

func jwtFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	token, _ := url.QueryUnescape(cookie.Value)
	if token == "" {
		return "", ErrEmptyToken
	}

	return token, nil
}

// Parse validates the token with the specialized secret,
// and store result to context
func Parse(tokenString string, secret string, ctx Context) (err error) {
	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return ctx.LoadFromMap(claims)
	}

	return
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

func GenerateToken(ctx Context) (tokenString string, err error) {
	// The token content.
	claims := jwt.MapClaims{}
	for k, v := range ctx.ToMap() {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(conf.Secret))

	return
}
