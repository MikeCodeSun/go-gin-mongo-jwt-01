package main

import (
	"example/go-gin-mongo-jwt/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err.Error())
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "home page"})
	})

	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	fmt.Println("hello")

	r.Run(":"+ port)
}