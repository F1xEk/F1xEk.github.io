package handlers

import (
	"charging-app/database"
	"charging-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

func Register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	query := `
		INSERT INTO users(name, email, password)
		VALUES($1, $2, $3)
	`

	_, err := database.DB.Exec(
		query,
		user.Name,
		user.Email,
		string(hashedPassword),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User creation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created",
	})
}

func Login(c *gin.Context) {
	var user models.User
	var dbUser models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}

	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email=$1
	`

	err := database.DB.QueryRow(
		query,
		user.Email,
	).Scan(
		&dbUser.ID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(dbUser.Password),
		[]byte(user.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}