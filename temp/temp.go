package main

import (
	"database/sql"
	"fmt"
	"strconv"
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
	fakeUser.FirstName = "Moe"
	fakeUser.LastName = "Jangda"

	fmt.Println(us.GetUser(1))

	err = us.UpdateUser(1, fakeUser)
	fmt.Println(err)
	defer db.Close()

	//fmt.Println(testQueryBuilder(fakeUser))

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

func testQueryBuilder(u *service.User) string {
	preparedUpdate := "UPDATE \"user\" SET "
	i := 0

	updots := make(map[string]int)
	values := make([]interface{}, 0)

	if u.FirstName != "" {
		i++
		updots["first_name"] = i
		values = append(values, u.FirstName)
	}

	if u.LastName != "" {
		i++
		updots["last_name"] = i
		values = append(values, u.LastName)
	}

	if u.Email != "" {
		i++
		updots["email"] = i
		values = append(values, u.Email)
	}

	if u.DateCreated != "" {
		i++
		updots["date_created"] = i
		values = append(values, u.DateCreated)
	}

	if u.Role != "" {
		i++
		updots["role"] = i
		values = append(values, u.Role)
	}

	if u.Password != "" {
		i++
		updots["password"] = i
		values = append(values, u.Password)
	}

	for k, v := range updots {
		preparedUpdate += k + "= $" + strconv.Itoa(v) + ", "
	}

	preparedUpdate += "WHERE id = $" + strconv.Itoa(len(values)+1)

	return preparedUpdate
}
