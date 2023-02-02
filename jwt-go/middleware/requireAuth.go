package middleware

import (
	"context"
	"fmt"
	"jwt-go/db/sqlc"
	"jwt-go/initializers"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")
	//get coockie off req
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	//decode/validate it

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)

		}
		//find tthe user with that token
		db := sqlc.New(initializers.DB)
		user, err := db.QueryGetUserById(context.Background(), int32(claims["sub"].(float64)))
		if err != nil || user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)

		}

		c.Set("user", user)
		//attach to req

		//continue

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
