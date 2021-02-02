package play

import (
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/middleware"
	"net/http"
)

//
//	This function sets up all routes for the play module.
//
func Init(router *gin.Engine) {
	authenticatedRoutes := router.Group("/play")

	authenticatedRoutes.Use(middleware.AuthenticateUserTokens(func (c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}))

	authenticatedRoutes.GET("/:key", GetPlayer)
}