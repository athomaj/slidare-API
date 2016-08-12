package authentication

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
	"models"
	"github.com/codegangsta/negroni"
	"mydb"
	"log"
)

var UserList = make([]models.UserModel, 0)

func doTokenExist(token *jwt.Token) (bool){
	return database.DoesTokenExist(token.Raw)
}

func RequireTokenAuthentication(userToken *string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		log.Println("sddsaahgfhgfs")
		authBackend := InitJWTAuthenticationBackend()

//		request.ParseFromRequest(req, request.OAuth2Extractor, keyLookupFunc);
		token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			} else {
				return authBackend.PublicKey, nil
			}
		})

		log.Println(err)
		log.Println(*userToken)
		if err == nil && token.Valid && doTokenExist(token) {
			*userToken = token.Raw
			next(rw, req)
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(err.Error()))
		}
	}
}
