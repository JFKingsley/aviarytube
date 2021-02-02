package browse

import (
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/middleware"
	"net/http"
)

//
//	This function sets up all routes for the browse module.
//
func Init(router *gin.Engine) {
	authenticatedRoutes := router.Group("/browse")

	authenticatedRoutes.Use(middleware.AuthenticateUserTokens(func (c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}))

	authenticatedRoutes.GET("/", GetBrowse)
}
