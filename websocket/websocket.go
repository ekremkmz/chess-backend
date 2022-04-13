package websocket

import (
	"chess-backend/websocket/model"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHandler(c *gin.Context) {
	req := c.Request

	cookie, err := req.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "access_token needed."})
		return
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "access_token is not valid."})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "access_token is not valid."})
		return
	}

	AddPlayer(&model.Player{
		Id:   claims["id"].(string),
		Conn: ws,
		Chan: make(chan []byte),
	})
}
