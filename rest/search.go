package rest

import (
	"chess-backend/restModels"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Search(c *gin.Context) {

	nick, ok := c.GetQuery("nick")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Bad request"})
		return
	}

	result := []restModels.User{}

	if mgm.Coll(&restModels.User{}).SimpleFindWithCtx(mgm.Ctx(), &result,
		bson.M{operator.Text: bson.M{"$search": nick}},
		&options.FindOptions{
			Sort:       bson.M{"score": bson.M{"$meta": "textScore"}},
			Projection: bson.M{"nick": 1, "score": bson.M{"$meta": "textScore"}, "_id": 0},
		}) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": true, "result": result})
}
