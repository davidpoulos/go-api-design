package service

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// User ...
type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName"  validate:"required"`
	Password    string `json:"password"  validate:"required"`
	Email       string `json:"email"     validate:"required"`
	DateCreated string `json:"dateCreated"`
	Role        string `json:"role"`
}

type UserService struct {
}

// EncryptPassword ...
func EncryptPassword(pass string) (string, error) {
	t, err := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if err != nil {
		return "", err
	}
	return string(t), nil
}

// UserStructLevelValidation ...
func UserStructLevelValidation(sl validator.StructLevel) {

	user := sl.Current().Interface().(User)

	// VALIDATE EVERYTHING HERE

	if user.Email == "David" || len(user.Password) == 0 {
		sl.ReportError(user.Email, "email", "email", "email", "")
		sl.ReportError(user.Password, "pasword", "Password", "pass", "")
	}

	// plus can do more, even with different tag than "fnameorlname"
}
