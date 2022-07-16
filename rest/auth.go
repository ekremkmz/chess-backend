package rest

import (
	"chess-backend/restModels"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// A function that register a user to mongo database
func Register(c *gin.Context) {
	var params struct {
		Nick     string `json:"nick"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := &restModels.User{}

	// Check if user already exists
	if mgm.Coll(user).FindOne(mgm.Ctx(), bson.M{operator.Or: bson.A{bson.M{"email": params.Email}, bson.M{"nick": params.Nick}}}).Decode(user) == nil {
		switch {
		case user.Email == params.Email:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email already exists"})
			return
		case user.Nick == params.Nick:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Nick already exists"})
			return
		}
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error hashing password"})
		return
	}

	// Create user
	user.Nick = params.Nick
	user.Email = params.Email
	user.Password = string(hashedPassword)

	// Save user
	if err := mgm.Coll(user).Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error saving user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Login(c *gin.Context) {
	var params struct {
		Nick     string `json:"nick"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	user := &restModels.User{}
	res := mgm.Coll(user).FindOne(mgm.Ctx(), bson.M{"nick": params.Nick})

	if err := res.Decode(user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	token, err := user.GenerateToken()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal server error"})
		return
	}

	httpOnly := false

	if os.Getenv("GO_ENV") == "development" {
		httpOnly = true
	}

	expire := os.Getenv("JWT_EXPIRE")

	exp, _ := strconv.ParseInt(expire, 10, 64)

	c.SetCookie("access_token", token, int(exp), "/", os.Getenv("DOMAIN"), false, httpOnly)

	c.JSON(http.StatusOK, gin.H{"success": true,
		"data": gin.H{"nick": user.Nick,
			"email":          user.Email,
			"friends":        user.Friends,
			"friendRequests": user.FriendRequests,
		},
	})
}

func Logout(c *gin.Context) {
	domain := os.Getenv("DOMAIN")

	httpOnly := false

	if os.Getenv("GO_ENV") == "development" {
		httpOnly = true
	}

	c.SetCookie("access_token", "", -1, "/", domain, false, httpOnly)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
