package entity

type User struct {
	ID          uint
	PhoneNumber string
	Name        string
	// Password always keep hash password.
	Password string
}
