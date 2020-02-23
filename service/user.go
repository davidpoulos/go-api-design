package service

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// User ...
type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"  validate:"required"`
	Password    string `json:"password"  validate:"required"`
	Email       string `json:"email"     validate:"required"`
	DateCreated string `json:"date_created"`
	Role        string `json:"role"`
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

	// TODO: VALIDATE EVERYTHING HERE -- Look at ways to do it in the metadata
	// IE. EMAIL 
	if user.Email == "David" || len(user.Password) == 0 {
		sl.ReportError(user.Email, "email", "email", "email", "")
		sl.ReportError(user.Password, "pasword", "Password", "pass", "")
	}

	// plus can do more, even with different tag than "fnameorlname"
}

// UserService ...
type UserService interface {
	InsertUser() error
	GetUser(id int) (*User, error)
}

// UserDB ...
type UserDB struct {
	db *sql.DB
}

// NewUserDB constructor
func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

// InsertUser ...
func (udb *UserDB) InsertUser(u User) error {

	userStmt := "INSERT INTO \"user\" (first_name, last_name, password, email, date_created, role) VALUES ($1,$2,$3,$4,$5,$6)"
	stmtIns, err := udb.db.Prepare(userStmt)
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(
		u.FirstName,
		u.LastName,
		u.Password,
		u.Email,
		u.DateCreated,
		u.Role)

	if err != nil {
		return err
	}

	return nil
}

// GetUser ...
func (udb *UserDB) GetUser(id int) (*User, error) {

	stmtOut, err := udb.db.Prepare("SELECT * FROM user WHERE id = $1 LIMIT 1")
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()

	var u User
	err = stmtOut.QueryRow(id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Password,
		&u.Email,
		&u.DateCreated,
		&u.Role)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
