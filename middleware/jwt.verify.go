package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var MySigningKey = []byte(os.Getenv("JWT_SECRET"))

func JwtVerify(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// load .env file
		err := godotenv.Load()
		if err != nil{
			log.Fatal("Error loading .env file")
		}
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf((" Invalid Signing Method"))
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}
				// verify iss claim
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}

				return MySigningKey, nil
			})
			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			if token.Valid {
				type context_key string
  ctx := context.WithValue(r.Context(), context_key("role_id"), token)
   next(w, r.WithContext(ctx))
			}

		} else {
			fmt.Fprintf(w, "Unauthorized!")
		}
	})
}
