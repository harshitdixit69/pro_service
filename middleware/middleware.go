package middleware

import (
	"context"
	"example/db"
	sqlc "example/db/sqlc"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie of req
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.HTML(200, "login.html", gin.H{})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.HTML(200, "login.html", gin.H{})
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.HTML(200, "login.html", gin.H{})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user with token sub
		user, _ := db.Query.Getusers(ctx, sqlc.GetusersParams{ID: int32(claims["sub"].(map[string]interface{})["id"].(float64)), IsUser: int64(claims["sub"].(map[string]interface{})["is_user"].(float64))})

		if user.ID == 0 {
			return
		}

		// Attach to req
		c.Set("user", user)

		// Continue
		c.Next()
	} else {
		c.HTML(200, "login.html", gin.H{})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
