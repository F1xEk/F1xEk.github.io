package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var statuses = []string{
	"Available",
	"Charging",
	"Offline",
}

func WSHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	defer conn.Close()

	ticker := time.NewTicker(3 * time.Second)

	defer ticker.Stop()

	energy := 0.0

	for {

		<-ticker.C

		energy += 0.12

		data := gin.H{
			"station_id": randomStation(),
			"status": randomStatus(),
			"power": rand.Intn(60) + 20,
			"energy": energy,
			"voltage": 400,
			"current": 32,
			"time": time.Now().Format("15:04:05"),
		}

		err := conn.WriteJSON(data)

		if err != nil {
			break
		}
	}
}

func randomStatus() string {
	return statuses[rand.Intn(len(statuses))]
}

func randomStation() int {
	return rand.Intn(3) + 1
}