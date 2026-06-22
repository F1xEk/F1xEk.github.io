package handlers

import (
	"charging-app/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {

	userID := c.MustGet("user_id")

	query := `
		SELECT id, name, email
		FROM users
		WHERE id=$1
	`

	var id int
	var name string
	var email string

	err := database.DB.QueryRow(
		query,
		userID,
	).Scan(
		&id,
		&name,
		&email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  name,
		"email": email,
	})
}