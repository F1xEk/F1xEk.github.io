package handlers

import (
	"charging-app/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StopSession(c *gin.Context) {

	type Request struct {
		SessionID int `json:"session_id"`
		Energy float64 `json:"energy"`
		Cost float64 `json:"cost"`
		Duration int `json:"duration"`
	}

	var req Request

	if err := c.BindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})

		return
	}

	_, err := database.DB.Exec(`
		UPDATE sessions
		SET
			status = 'completed',
			energy = $1,
			cost = $2,
			duration = $3
		WHERE id = $4
	`,
		req.Energy,
		req.Cost,
		req.Duration,
		req.SessionID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to stop session",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session completed",
	})
}