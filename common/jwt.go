package common

import (
	"github.com/dgrijalva/jwt-go"
	"jack/model"
	"time"
)

var jwtKey = []byte("Zhejiang")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User)(string,error)  {
	expirationTime := time.Now().Add(7 * 24 *time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "jack ping",
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString,nil
}

func ParseToken(tokenString string) (*jwt.Token,*Claims,error) {
	claims := &Claims{}

	token,err := jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey,nil

	})
	return token,claims,err
}