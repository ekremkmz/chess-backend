package main

import (
	"chess-backend/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRoutes() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	r := gin.Default()

	r.LoadHTMLFiles("index.html")

	r.GET("/", rootHandler)

	r.GET("/ws", websocket.WebsocketHandler)

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	setupRoutes()

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
