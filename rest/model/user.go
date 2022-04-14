package model

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Nick             string `bson:"nick" json:"nick"`
	Email            string `bson:"email" json:"email"`
	Password         string `bson:"password" json:"password"`
}

func (u *User) GenerateToken() (string, error) {
	jwt.New(jwt.SigningMethodHS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nick": u.Nick,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
