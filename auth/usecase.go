package auth

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeri06/majoo/exception"
	"github.com/jeri06/majoo/model"
	"github.com/jeri06/majoo/response"
)

var jwtKey = []byte("123456")

type Usecase interface {
	Auth(ctx context.Context, auth model.Auth) (resp response.Response)
}

type authUsecase struct {
	repository Repository
}

func NewAuthUsecase(rp Repository) Usecase {
	return &authUsecase{
		repository: rp,
	}
}

func (u authUsecase) Auth(ctx context.Context, auth model.Auth) (resp response.Response) {
	user, err := u.repository.FindByUsername(ctx, auth.UserName)
	if err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	outlets, err := u.repository.FindOutletByMerchantId(ctx, user.MerchandId)
	if err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	user.OutletId = outlets

	if auth.Password != user.Password {
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, "")
	}
	expirationTime := time.Now().Add(100 * time.Minute)

	claims := &model.Claims{
		Authorization: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, "")

	}
	resp = response.NewSuccessResponse(tokenString, response.StatOK, "")
	return
}
