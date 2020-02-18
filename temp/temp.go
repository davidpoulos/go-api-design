package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/davidpoulos/hackin/service"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := GetMySQLDB()

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

}

// InsertUser ...
func InsertUser(db *sql.DB, u service.User) error {

	fmt.Printf("User: %v", u)

	userStmt := "INSERT INTO User (firstName, lastName, password, email, dateCreated, role) VALUES(?,?,?,?,?,?)"
	stmtIns, err := db.Prepare(userStmt)
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.FirstName, u.LastName, u.Password, u.Email, u.DateCreated, u.Role)

	if err != nil {
		return err
	}

	return nil
}

// EncryptPassword ...
func EncryptPassword(pass string) (string, error) {
	t, err := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if err != nil {
		return "", err
	}
	return string(t), nil
}

// GetMockUser ...
func GetMockUser() (u *service.User) {
	t := time.Now().UTC()
	p, _ := EncryptPassword("pepperoni")
	return &service.User{
		FirstName: "David",
		LastName:  "Poulos",
		DateCreated: fmt.Sprintf("%d-%02d-%02d",
			t.Year(), t.Month(), t.Day()),
		Email:    "david@datafiniti.co",
		Role:     "admin",
		Password: p,
	}

}

// GetMySQLDB ...
func GetMySQLDB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "zues"
	dbName := "testdb"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
