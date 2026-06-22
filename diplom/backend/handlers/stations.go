package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStations(c *gin.Context) {

	stations := []gin.H{
		{
			"id": 1,
			"title": "E-Prom",
			"address": "Красносельская, 51к4",
			"latitude": 55.817449,
			"longitude": 49.096843,
			"status": "Свободна",
			"connector": "CHAdeMO",
			"tariff": "15 р/кВтч",
		},
		{
			"id": 2,
			"title": "Evolute",
			"address": "Московская, 20",
			"latitude": 55.784368,
			"longitude": 49.187622,
			"status": "Занята",
			"connector": "CCS",
			"tariff": "15 р/кВтч",
		},
		{
			"id": 3,
			"title": "Pandora",
			"address": "Аделя Кутуя, 83а",
			"latitude": 55.790866,
			"longitude": 49.105300,
			"status": "Свободна",
			"connector": "CCS",
			"tariff": "15 р/кВтч",
		},
	}

	c.JSON(http.StatusOK, stations)
}