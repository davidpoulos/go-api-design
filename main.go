package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/davidpoulos/hackin/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	server := gin.Default()

	_, err := GetMySQLDB()

	if err != nil {
		panic("Un-able to connect to MySQL")
	}

	fmt.Println("Connected to MYSQL")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(service.UserStructLevelValidation, service.User{})
	}

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.POST("/user", func(c *gin.Context) {

		var u service.User

		if err := c.ShouldBindJSON(&u); err != nil {
			errs := make([]string, 0)
			for _, fieldErr := range err.(validator.ValidationErrors) {
				errs = append(errs, fieldErr.Field())
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return // exit on first error
		}

		fmt.Println(c.Get("user"))

		fmt.Printf("%v", u)

		c.Writer.WriteHeader(http.StatusCreated)

	})

	server.Run()

}

// AuthMiddleware ... TODO
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set("user", gin.H{"id": "234234"})
	}
}

// GetMySQLDB ...
func GetMySQLDB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "zeus"
	dbName := "tesdb"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}
