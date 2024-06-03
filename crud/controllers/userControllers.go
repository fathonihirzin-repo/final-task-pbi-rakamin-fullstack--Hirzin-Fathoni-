package controllers

import (
	"crud/initializers"
	"crud/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// Get email/pass

	var body struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	// Create user

	user := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//Respond
	c.JSON(http.StatusOK, gin.H{"": "user created"})

}

func Login(c *gin.Context) {

	//Get email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	//Look up request
	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}

	// Compare sent in password with save user pass

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}

	// Generate jwt token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
	}

	//send it back to cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		// "token": tokenString,
	})

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		// "messages": user,
		"hallo": user.(models.User).Username,
	})
}

func UserIndex(c *gin.Context) {
	//Get all users
	var users []models.User
	initializers.DB.Find(&users)

	//Respond
	c.JSON(200, gin.H{
		"Users": users,
	})
}

func GetUser(c *gin.Context) {
	//Get id from url
	id := c.Param("id")

	//Get all users
	var user models.User
	initializers.DB.Find(&user, id)

	//Respond
	c.JSON(200, gin.H{
		"Users": user,
	})
}

func UpdateUser(c *gin.Context) {
	//Get id from url
	id := c.Param("id")

	//Get data off req body
	var body struct {
		Username string
		Email    string
		Password string
	}

	c.Bind(&body)

	//Find the user were updating
	var user models.User
	initializers.DB.Find(&user, id)

	//Update
	initializers.DB.Model(&user).Updates(models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})

	//Respond
	c.JSON(200, gin.H{
		"Users": user,
	})

}

func DeleteUser(c *gin.Context) {
	//Get id from url
	id := c.Param("id")

	//Delete the posts
	initializers.DB.Delete(&models.User{}, id)

	//Respond
	c.JSON(200, gin.H{
		"messages": "user deleted",
	})
}

// type PhotoRequest struct {
// 	ID int `json:id`
// }
