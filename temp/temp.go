package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/davidpoulos/hackin/service"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func main() {
	db, err := GetPostGresDB()

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	us := service.NewUserDB(db)

	fakeUser := GetMockUser()
	
	e := us.InsertUser(*fakeUser)
	fmt.Println(e)

	//fmt.Println(us.GetUser(2))

	defer db.Close()

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
func GetPostGresDB() (db *sql.DB, err error) {
	user := "root"
	pass := "zues"
	dbname := "testdb"
	host := "localhost"
	port := 5432

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, pass, dbname)
	
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
