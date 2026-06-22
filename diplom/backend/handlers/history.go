package handlers

import (
	"charging-app/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSessions(c *gin.Context) {

	userID := c.MustGet("user_id")

	rows, err := database.DB.Query(`
		SELECT
			id,
			station_id,
			status,
			energy,
			cost,
			duration,
			started_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
	`, userID)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load sessions",
		})

		return
	}

	defer rows.Close()

	var sessions []gin.H

	for rows.Next() {

		var (
			id int
			stationID int
			status string
			energy float64
			cost float64
			duration int
			startedAt string
		)

		rows.Scan(
			&id,
			&stationID,
			&status,
			&energy,
			&cost,
			&duration,
			&startedAt,
		)

		sessions = append(sessions, gin.H{
			"id": id,
			"station_id": stationID,
			"status": status,
			"energy": energy,
			"cost": cost,
			"duration": duration,
			"started_at": startedAt,
		})
	}

	c.JSON(http.StatusOK, sessions)
}