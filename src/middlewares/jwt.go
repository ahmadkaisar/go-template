package middlewares

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"

	handler "../handlers"
)

type JWT struct {}

var mysigningkey = []byte("SomethingInteresting123456")

func (j JWT) Generate(user_id, role_id int, user_name, role_name string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["user_name"] = user_name
	claims["role_id"] = role_id
	claims["role_name"] = role_name
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	tokenstring, err := token.SignedString(mysigningkey)

	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

func (j JWT) Claim(t string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return mysigningkey, nil
	})

	if err != nil {
		return nil, err
	}

	// for key, val := range claims {
	// 	fmt.Printf("-Key: %v, value: %v\n", key, val)
	// }
	return claims, nil
}

func (j JWT) Validate(endpoint http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Method == "OPTIONS" {
			handler.Response(w, r, 200, "allowed")
			return
		} else if r.Header.Get("Authorization") != "" {
			bearer := strings.Split(r.Header.Get("Authorization"), " ")
			if len(bearer) != 2 {
				handler.Response(w, r, 500, "Authorization: Bearer <token>")
				return
			}

			token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error when parsing token")
				}

				return mysigningkey, nil
			})
			if err != nil {
				handler.Response(w, r, 500, "there was an error when parsing token")
				return
			} else if token.Valid {
				endpoint(w, r)
				return
			} else {
				handler.Response(w, r, 500, "token not valid")
				return
			}
		}
		cookies, err := r.Cookie("Authorization")
		if err != nil {
			handler.Response(w, r, 401, "not authorized")
			return
		} else {
			bearer := strings.Split(cookies.Value, " ")
			if len(bearer) != 2 {
				handler.Response(w, r, 500, "Authorization: Bearer <token>")
				return
			}

			token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error when parsing token")
				}

				return mysigningkey, nil
			})
			if err != nil {
				handler.Response(w, r, 500, "there was an error when parsing token")
				return
			} else if token.Valid {
				endpoint(w, r)
				return
			} else {
				handler.Response(w, r, 500, "token not valid")
				return
			}
		}
		handler.Response(w, r, 401, "not authorized")
	})
}
