package rest

import (
	"chess-backend/rest/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func Profile(c *gin.Context) {
	token, _ := c.Get("token")

	nick := token.(jwt.MapClaims)["nick"].(string)

	user := &model.User{}

	if mgm.Coll(user).FindOne(mgm.Ctx(), bson.M{"nick": nick}).Decode(user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User not found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"nick": user.Nick, "email": user.Email}})
}
