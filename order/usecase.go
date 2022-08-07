package order

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/jeri06/majoo/model"
	"github.com/jeri06/majoo/response"
)

type Usecase interface {
	ReportOmzetMerchant(ctx context.Context, param model.DailyReportMerchantParam) (resp response.Response)
	ReportOmzetoutlets(ctx context.Context, param model.DailyReportOutletParam) (resp response.Response)
}

type orderUsecase struct {
	repository Repository
}

func NewOrderUsecase(rp Repository) Usecase {
	return &orderUsecase{
		repository: rp,
	}
}

func (u orderUsecase) ReportOmzetMerchant(ctx context.Context, param model.DailyReportMerchantParam) (resp response.Response) {

	userAccess, ok := ctx.Value(model.AdminKey{}).(model.Claims)
	if !ok {
		return response.NewErrorResponse(nil, http.StatusUnauthorized, nil, response.StatUnauthorized, "")

	}

	offset := (param.Page - 1) * param.Size
	dailyReport, err := u.repository.ReportOmzetByMerchant(ctx, param, int(userAccess.MerchandId), int64(offset), int64(param.Size))
	if err != nil {
		fmt.Println(err.Error())
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	totalData, err := u.repository.ReportOmzetByMerchantTotalData(ctx, param, int(userAccess.MerchandId))
	if err != nil {
		fmt.Println(err.Error())
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	totalDataOnPage := len(dailyReport)
	totalPage := int(math.Ceil(float64(totalData) / float64(param.Size)))
	meta := Meta{
		Page:            int(param.Page),
		TotalPage:       totalPage,
		TotalData:       totalData,
		TotalDataOnPage: totalDataOnPage,
	}

	resp = response.NewSuccessResponseWithMeta(dailyReport, meta, response.StatOK, "")
	return
}

func (u orderUsecase) ReportOmzetoutlets(ctx context.Context, param model.DailyReportOutletParam) (resp response.Response) {
	var eligibleOutlet *int64

	userAccess, ok := ctx.Value(model.AdminKey{}).(model.Claims)
	if !ok {
		fmt.Println("ok")
		return response.NewErrorResponse(nil, http.StatusUnauthorized, nil, response.StatUnauthorized, "")

	}

	for _, v := range userAccess.OutletId {
		if v == param.OutletId {
			eligibleOutlet = &v
		}
	}
	if eligibleOutlet == nil {
		return response.NewErrorResponse(nil, http.StatusUnauthorized, nil, response.StatUnauthorized, "")
	}
	offset := (param.Page - 1) * param.Size
	dailyReport, err := u.repository.ReportOmzetOutlets(ctx, param, int(userAccess.MerchandId), *eligibleOutlet, int64(offset), int64(param.Size))
	if err != nil {
		fmt.Println(err.Error())
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	totalData, err := u.repository.ReportOmzetOutletsTotalData(ctx, param, int(userAccess.MerchandId))
	if err != nil {
		fmt.Println(err.Error())
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	totalDataOnPage := len(dailyReport)
	totalPage := int(math.Ceil(float64(totalData) / float64(param.Size)))
	meta := Meta{
		Page:            int(param.Page),
		TotalPage:       totalPage,
		TotalData:       totalData,
		TotalDataOnPage: totalDataOnPage,
	}

	resp = response.NewSuccessResponseWithMeta(dailyReport, meta, response.StatOK, "")
	return
}
