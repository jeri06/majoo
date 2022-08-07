package order

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jeri06/majoo/middleware"
	"github.com/jeri06/majoo/model"
	"github.com/jeri06/majoo/response"
)

type HTTPHandler struct {
	Usecase Usecase
}

func NewHTTPHandlerOrder(router *mux.Router, usecase Usecase, session middleware.RouteMiddleware) {
	handle := &HTTPHandler{
		Usecase: usecase,
	}

	router.HandleFunc("/order-service/v1/DailyReportMerchant", session.Verify(handle.ReportMerchant)).Methods(http.MethodGet)
	router.HandleFunc("/order-service/v1/DailyReportOutlets", session.Verify(handle.ReportOutlet)).Methods(http.MethodGet)
}

func (handler *HTTPHandler) ReportMerchant(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	ctx := r.Context()

	params := model.DailyReportMerchantParam{}
	params.StartDate = queryString.Get("startDate")
	params.EndDate = queryString.Get("endDate")
	params.Page, _ = strconv.ParseInt(queryString.Get("page"), 10, 64)
	params.Size, _ = strconv.ParseInt(queryString.Get("size"), 10, 64)
	resp := handler.Usecase.ReportOmzetMerchant(ctx, params)
	response.JSON(w, resp)

}

func (handler *HTTPHandler) ReportOutlet(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	ctx := r.Context()

	params := model.DailyReportOutletParam{}
	params.StartDate = queryString.Get("startDate")
	params.EndDate = queryString.Get("endDate")
	params.Page, _ = strconv.ParseInt(queryString.Get("page"), 10, 64)
	params.Size, _ = strconv.ParseInt(queryString.Get("size"), 10, 64)
	params.OutletId, _ = strconv.ParseInt(queryString.Get("outletId"), 10, 64)

	resp := handler.Usecase.ReportOmzetoutlets(ctx, params)
	response.JSON(w, resp)

}
