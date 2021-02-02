package login

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	"github.com/ulule/limiter/drivers/store/memory"
	"net/http"
	"time"
	"github.com/gin-contrib/sessions"
)

var store = memory.NewStore()

//
//	This function sets up all routes for the login module.
//
func Init(router *gin.Engine) {
	baseRoutes := router.Group("/login")

	// Setup a rate limit to avoid user enumeration
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	// // Then, create the limiter instance which takes the store and the rate as arguments.
	// // Now, you can give this instance to any supported middleware.
	// limitedRoute := router.Group("/ui/v1")
	limitedMiddleware := mgin.NewMiddleware(limiter.New(store, rate), mgin.WithLimitReachedHandler(func(c *gin.Context) {
		session := sessions.Default(c)
		session.AddFlash("You have exceeded your login attempts. Please wait a minute and try again.")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}))

	baseRoutes.Use(limitedMiddleware)

	baseRoutes.GET("/", GetLogin)
	baseRoutes.GET("/invalidate", GetLogout)
	baseRoutes.GET("/oauth", OAuthRedirect)
}
