package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CurrentUser(ctx *gin.Context) {
	tokenstring, err := ctx.Cookie("auth")
	if err == nil {
		ctx.Set("currentUserToken", tokenstring)
	} else {
		ctx.Set("currentUserToken", nil)
	}
	ctx.Next()
}