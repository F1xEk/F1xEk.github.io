package handlers

import (
	"charging-app/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(c *gin.Context) {

	userID := c.MustGet("user_id")

	var (
		totalSessions int
		totalEnergy float64
		totalCost float64
		avgDuration float64
	)

	err := database.DB.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(energy), 0),
			COALESCE(SUM(cost), 0),
			COALESCE(AVG(duration), 0)
		FROM sessions
		WHERE user_id = $1
	`,
		userID,
	).Scan(
		&totalSessions,
		&totalEnergy,
		&totalCost,
		&avgDuration,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Analytics error",
		})

		return
	}

	rows, err := database.DB.Query(`
		SELECT
			started_at,
			energy
		FROM sessions
		WHERE user_id = $1
		ORDER BY started_at
	`,
		userID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Chart data error",
		})

		return
	}

	defer rows.Close()

	var chartData []gin.H

	for rows.Next() {

		var (
			date string
			energy float64
		)

		rows.Scan(
			&date,
			&energy,
		)

		chartData = append(chartData, gin.H{
			"date": date,
			"energy": energy,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_sessions": totalSessions,
		"total_energy": totalEnergy,
		"total_cost": totalCost,
		"avg_duration": avgDuration,
		"chart": chartData,
	})
}