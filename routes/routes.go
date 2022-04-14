package routes

import (
	"example/go-gin-mongo-jwt/controllers"
	"example/go-gin-mongo-jwt/utils"

	"github.com/gin-gonic/gin"
)


func UserRoutes(r *gin.Engine) {
	r.POST("/users/register", controllers.Register())
	r.POST("/users/login", controllers.Login())
}

func AuthRoutes(r *gin.Engine) {
	r.GET("/users", utils.Auth(), controllers.GetAllUsers())
	r.GET("/users/:id",utils.Auth(), controllers.GetUserById())
}