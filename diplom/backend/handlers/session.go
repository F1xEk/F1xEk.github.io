package handlers

import (
	"charging-app/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartSession(c *gin.Context) {

	userID := c.MustGet("user_id")

	type Request struct {
		StationID int `json:"station_id"`
	}

	var req Request

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	query := `
		INSERT INTO sessions (
			user_id,
			station_id,
			status
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var sessionID int

	err := database.DB.QueryRow(
		query,
		userID,
		req.StationID,
		"charging",
	).Scan(&sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Charging started",
		"session_id": sessionID,
	})
}