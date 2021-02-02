package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticateUserTokens(errorFunction func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		v := session.Get("user-id")

		if v == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		}
		c.Next()
	}
}