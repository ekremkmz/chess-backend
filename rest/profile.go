package rest

import (
	"chess-backend/restModels"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Profile(c *gin.Context) {
	token, _ := c.Get("token")

	nick := token.(jwt.MapClaims)["nick"].(string)

	user := &restModels.User{}

	if mgm.Coll(user).FindOne(mgm.Ctx(), bson.M{"nick": nick}, &options.FindOneOptions{
		Projection: bson.M{"password": 0},
	}).Decode(user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User not found"})
		return
	}
	email := user.Email
	friends := user.Friends
	fRequests := user.FriendRequests

	c.JSON(http.StatusOK,
		gin.H{"success": true,
			"data": gin.H{"nick": nick,
				"email":          email,
				"friends":        friends,
				"friendRequests": fRequests,
			},
		})
}
