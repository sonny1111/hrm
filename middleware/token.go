package middleware

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

//Authentication function will call this middleware for generating jwt token
func GenerateJWT(username, roleId string) (string, error) {
	// load .env file
	err := godotenv.Load()
if err != nil{
	log.Fatal("Error loading .env file")
}
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))
	Token := jwt.New(jwt.SigningMethodHS256)
	claims := Token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = username
	claims["roleId"] = roleId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := Token.SignedString(mySigningKey)
	if err != nil {
	log.Fatal("Unauthorized!")
	}
	return tokenString, nil
}

