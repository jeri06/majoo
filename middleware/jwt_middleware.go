package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeri06/majoo/exception"
	"github.com/jeri06/majoo/model"
	"github.com/jeri06/majoo/response"
)

var (
	header = "Authorization"
	jwtKey = []byte("123456")
)

type Session struct{}

func NewJwtMiddleware() RouteMiddleware {
	return &Session{}
}

func (a Session) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get(header)

		if authHeader == "" {
			a.respondUnauthorized(w, jwt.ErrInvalidKey.Error())
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			a.respondUnauthorized(w, "not valid")
			return
		}
		claims := model.Claims{}
		tkn, err := jwt.ParseWithClaims(bearerToken[1], &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			a.respondUnauthorized(w, err.Error())
			return
		}
		if !tkn.Valid {
			a.respondUnauthorized(w, jwt.ErrSignatureInvalid.Error())
			return
		}
		context := context.WithValue(ctx, model.AdminKey{}, claims)
		r = r.WithContext(context)

		next(w, r)
	})

}

func (a Session) respondUnauthorized(w http.ResponseWriter, message string) {
	resp := response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, message)
	response.JSON(w, resp)
}
