package main

import (
	"context"
	"fmt"
	"gosome/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.Default()

	router.Use(dbMiddleware(*conn))

	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UserRegister)
		usersGroup.POST("login", routes.UserLogin)
	}

	router.Run(":3001")

}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:@localhost:5432/someap")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err.Error())
	}

	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandleFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}
