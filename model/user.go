package model

type Auth struct {
	// Id       string `json:"Id"`
	// Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Authorization struct {
	Name         string  `json:"name"`
	UserName     string  `json:"user_name"`
	Password     string  `json:"password"`
	MerchandId   int64   `json:"merchantId"`
	MerchandName string  `json:"merchant_name"`
	OutletId     []int64 `json:"outletId"`
}

type (
	Area struct {
		ID        int64  `gorm:"column:id;primaryKey;"`
		AreaValue int64  `gorm:"column:area_value"`
		AreaType  string `gorm:"column:type"`
	}
)
