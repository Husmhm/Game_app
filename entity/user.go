package entity

type User struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	// Password always keep hash password.
	Password string
	Role     Role
}
