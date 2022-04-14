package controllers

import (
	"context"
	"example/go-gin-mongo-jwt/database"
	"example/go-gin-mongo-jwt/models"
	"example/go-gin-mongo-jwt/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
var client = database.Client
var collection = database.ConnectCollection(client, "user")


func Register() gin.HandlerFunc{
  return func(c *gin.Context) {
		var user models.User
		
		// get data input json
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		}
		// validate user input
	 err :=	validate.Struct(user)
	 if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return	
	 }
	 // 
	 ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
	//  check user exist ?
	 count, err := collection.CountDocuments(ctx, bson.M{"name": user.Name})
	 defer cancel()
	 if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return	
	 }
	 if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user exist"})
		return	
	 } 
	 	//  bcrypt hash password
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		 }

		user.Password = string(hashPassword)
	  user.ID = primitive.NewObjectID()

	  
		 res, err := collection.InsertOne(ctx, user)
		 defer cancel()
		 if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		 }

		 c.JSON(http.StatusOK, res)
  }
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User
// bind user
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		}
// context 
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
		// find user exist ?
		err := collection.FindOne(ctx, bson.M{"name": user.Name}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		 }
		 if foundUser.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "no user"})
			return	
		 }
	// compare user password
	  errH := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if errH != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": errH.Error()})
			return	
		 }
		 // generate token
		 token := utils.GenerateJwt(foundUser.Name)
		 c.JSON(http.StatusOK, token)
	}
}

func GetAllUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		var users []models.User
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)

		cur, err := collection.Find(ctx, bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		}

		for cur.Next(ctx) {
			var user models.User
			if err := cur.Decode(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
			}
			users = append(users, user)
		}
		defer cancel()
		c.JSON(http.StatusOK, users)
	}
}

func GetUserById() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User
		id := c.Param("id")

		userid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)

		res := collection.FindOne(ctx, bson.M{"_id": userid})
		defer cancel()
		// if res == nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"msg": "no user"})
		// 	return	
		// }
		
		if err := res.Decode(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return	
		}

		c.JSON(http.StatusOK, user)

	}
}