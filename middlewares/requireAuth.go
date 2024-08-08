package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(ctx *gin.Context) {
	cookie, exists := ctx.Get("currentUserToken")
	if !exists || cookie == nil {
		fmt.Println("no user found")
		ctx.JSON(http.StatusOK, gin.H{})
		ctx.Abort()
		return
	}
	tokenstring, ok := cookie.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		ctx.Abort()
		return
	}
	if err := validateToken(tokenstring, ctx); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return
	}


	ctx.Next()
}

func validateToken(tokenString string, ctx *gin.Context) (error) {
	// Retrieve the secret key from environment variables or configuration
	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		return errors.New("missing secret key")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return errors.New("token is expired")
		}
	}

	if err != nil {
		return err
	}

	// Check if the token is valid
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}