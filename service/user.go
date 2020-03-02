package service

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// User ...
type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"  validate:"required"`
	Password    string `json:"password"  validate:"required"`
	Email       string `json:"email"     validate:"required"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	Role        string `json:"role"`
}

// NewUser ...
func NewUser() *User {
	t := time.Now()
	return &User{
		DateUpdated: fmt.Sprintf("%d-%02d-%02d",
			t.Year(), t.Month(), t.Day()),
	}
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

// UpdateUser ...
func (udb *UserDB) UpdateUser(id int, u *User) error {

	preparedUpdate := "UPDATE \"user\" SET "

	fields := make(map[string]string, 0)

	if u.FirstName != "" {
		fields["first_name"] = u.FirstName
	}

	if u.LastName != "" {
		fields["last_name"] = u.LastName
	}

	if u.Email != "" {
		fields["email"] = u.Email
	}

	if u.DateCreated != "" {
		fields["date_created"] = u.DateCreated
	}

	if u.Role != "" {
		fields["role"] = u.Role
	}

	if u.Password != "" {
		fields["password"] = u.Password
	}

	i := 1
	values := make([]interface{}, 0)
	for k, v := range fields {
		preparedUpdate = preparedUpdate + k + "= $" + strconv.Itoa(i) + ","
		values = append(values, v)
		i++
	}

	preparedUpdate = strings.TrimRight(preparedUpdate, ", ")
	preparedUpdate += "WHERE id = $" + strconv.Itoa(len(values)+1)

	values = append(values, strconv.Itoa(id))

	fmt.Println(preparedUpdate)
	_, err := udb.db.Query(preparedUpdate, values...)

	if err != nil {
		return err
	}

	return nil
}

// GetUser ...
func (udb *UserDB) GetUser(id int) (*User, error) {

	stmtOut, err := udb.db.Prepare("SELECT * FROM \"user\" WHERE id = $1 LIMIT 1")
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
