package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeri06/majoo/middleware"
	"github.com/jeri06/majoo/model"
	"github.com/jeri06/majoo/response"
)

type HTTPHandler struct {
	Usecase Usecase
}

func NewHTTPHandlerAuth(router *mux.Router, usecase Usecase, session middleware.RouteMiddleware) {
	handle := &HTTPHandler{
		Usecase: usecase,
	}

	router.HandleFunc("/order-service/v1/auth", handle.Auth).Methods(http.MethodPost)

}

func (handler HTTPHandler) Auth(w http.ResponseWriter, r *http.Request) {

	var resp response.Response
	var auth model.Auth
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = handler.Usecase.Auth(ctx, auth)
	response.JSON(w, resp)
	return
}
