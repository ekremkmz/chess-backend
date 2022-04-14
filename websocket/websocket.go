package websocket

import (
	"chess-backend/websocket/model"
	"net/http"

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

	ws, err := upgrader.Upgrade(c.Writer, req, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "access_token is not valid."})
		return
	}

	token, _ := c.Get("token")

	nick := token.(jwt.MapClaims)["nick"].(string)

	AddPlayer(&model.Player{
		Nick: nick,
		Conn: ws,
		Chan: make(chan []byte),
	})
}
