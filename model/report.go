package model

type DailyReportMerchant struct {
	Date         string  `json:"date"`
	MerchantName string  `json:"merchantName"`
	Omzet        float32 `json:"omzet"`
}

type DailyReportOutlet struct {
	Date         string  `json:"date"`
	MerchantName string  `json:"merchantName"`
	OutletName   string  `json:"outlet_name"`
	Omzet        float32 `json:"omzet"`
}

type DailyReportMerchantParam struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Page      int64  `json:"page"`
	Size      int64  `json:"size"`
}

type DailyReportOutletParam struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Page      int64  `json:"page"`
	Size      int64  `json:"size"`
	OutletId  int64  `json:"outletId"`
}
