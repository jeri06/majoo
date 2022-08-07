package entity

type User struct {
	Id       string `json:"Id"`
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
