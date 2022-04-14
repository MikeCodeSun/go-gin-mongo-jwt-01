package utils

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Claims struct {
	Name string
	jwt.StandardClaims
}

func GenerateJwt(name string) string{
 if err := godotenv.Load(); err != nil {
	 log.Fatal(err)
 }
 secret := os.Getenv("SECRET")

  claims := Claims{
		Name: name,
		// standard claim jwt
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(24 * time.Hour).Unix(),
		},
	}

	token, err:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

return token
}

func Auth() gin.HandlerFunc{
	return func(c *gin.Context){
		var claims Claims
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		secret := os.Getenv("SECRET")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "no header auth"})
			c.Abort()
			return
		}
		token := strings.Split(authHeader, " ")[1]
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "no token auth"})
			c.Abort()
			return
		}

		 t,err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret),nil
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			c.Abort()
			return
		}

		if t.Valid != true {
			log.Fatalln(t)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "auth not right"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}