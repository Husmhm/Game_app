package param

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	User   UserInfo `json:"mysqluser"`
	Tokens Tokens   `json:"tokens"`
}
