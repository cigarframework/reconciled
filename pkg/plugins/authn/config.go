package authn

import (
	"errors"
	"strings"
	"time"

	"github.com/cigarframework/reconciled/pkg/api"
	jwt "github.com/dgrijalva/jwt-go"
)

type Config struct {
	CacheDuration time.Duration `yaml:"cacheDuration"`
	JWT           *JWT          `yaml:"jwt"`
}

type JWT struct {
	Key string `yaml:"key"`
}

type AuthMethod interface {
	Auth(token string) (*api.User, error)
}

type AuthMethodFunction func(token string) (*api.User, error)

func (a AuthMethodFunction) Auth(token string) (*api.User, error) {
	return a(token)
}

type jwtAuth struct {
	key interface{}
}

func (a *jwtAuth) Auth(token string) (*api.User, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return a.key, nil
	})
	if err != nil {
		return nil, err
	}
	parts := strings.Split(claims.Subject, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid jwt claims")
	}
	return &api.User{
		Name:  parts[0],
		Group: parts[1],
	}, nil
}
