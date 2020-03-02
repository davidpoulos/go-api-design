package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidpoulos/hackin/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {

	server := gin.Default()

	db, err := GetPostGresDB()

	if err != nil {
		panic("Un-able to connect to Postgres")
	}

	us := service.NewUserDB(db)

	log.Println("==> Connected to Postgres")

	// Register Validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(service.UserStructLevelValidation, service.User{})
	}

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "green",
		})
	})

	server.POST("/auth", func(c *gin.Context) {
		//Find user w/

		bitties, _ := c.GetRawData()

		fmt.Println(string(bitties))

	})

	server.POST("/user", func(c *gin.Context) {

		u := service.NewUser()
		t := time.Now()

		u.DateCreated = fmt.Sprintf("%d-%02d-%02d",
			t.Year(), t.Month(), t.Day())

		if err := c.ShouldBindJSON(u); err != nil {
			errs := make([]string, 0)
			for _, fieldErr := range err.(validator.ValidationErrors) {
				errs = append(errs, fieldErr.Field())
			}
			c.JSON(http.StatusConflict, gin.H{"errors": errs})
			return
		}

		u.Password, _ = service.EncryptPassword(u.Password)
		err = us.InsertUser(*u)
		if err, ok := err.(*pq.Error); ok {

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": formalizePostgresErrorMessage(err),
			})
		} else {
			c.Writer.WriteHeader(http.StatusCreated)
		}

	})

	server.PUT("/user", func(c *gin.Context) {

	})

	server.Run()

}

func formalizePostgresErrorMessage(p *pq.Error) string {
	errMsg := "Error Validating Payload"
	if p.Code == "23505" {
		errMsg = "Email already in use"
	}

	return errMsg
}

// AuthMiddleware ... TODO
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set("user", gin.H{"id": "234234"})
	}
}

// GetPostGresDB ...
func GetPostGresDB() (db *sql.DB, err error) {
	user := "root"
	pass := "zues"
	dbname := "testdb"
	host := "postgres"
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
