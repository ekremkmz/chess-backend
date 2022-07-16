package main

import (
	"chess-backend/rest"
	"chess-backend/websocket"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := setupMgm(); err != nil {
		log.Fatal(err)
	}

	r := setupRoutes()

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func setupMgm() error {
	host := os.Getenv("DB_HOST")
	auth := options.Credential{Username: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASS")}
	clientOpt := options.Client().ApplyURI(host).SetAuth(auth)

	if err := mgm.SetDefaultConfig(nil, os.Getenv("DB_NAME"), clientOpt); err != nil {
		return err
	}

	return nil
}

func setupRoutes() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLFiles("index.html")

	r.GET("/socket_test", socketTester)

	r.GET("/ws", AuthNeeded, websocket.WebsocketHandler)

	r.GET("/profile", AuthNeeded, rest.Profile)

	auth := r.Group("/auth")

	{
		auth.POST("/register", rest.Register)
		auth.POST("/login", rest.Login)
		auth.GET("/logout", AuthNeeded, rest.Logout)
	}

	r.GET("/search", AuthNeeded, rest.Search)

	return r
}

func socketTester(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
