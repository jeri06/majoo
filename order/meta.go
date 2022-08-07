package order

type Meta struct {
	Page            int `json:"page"`
	TotalPage       int `json:"totalPage"`
	TotalData       int `json:"totalData"`
	TotalDataOnPage int `json:"totalDataOnPage"`
}
