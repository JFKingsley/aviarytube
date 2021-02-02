package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/configuration"
	"github.com/jfkingsley/aviarytube/views/browse"
	"github.com/jfkingsley/aviarytube/views/login"
	"github.com/jfkingsley/aviarytube/views/play"
	"github.com/yargevad/filepathx"
	"go.uber.org/zap"
)

//
//  This adapter proxies any API Gateway request objects to Gin for standard routing.
//
var ginLambda *ginadapter.GinLambda

func init() {
	// Setup a new Zap logger.
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Setup our Gin instance.
	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "state": "Success",
        })
    })

    router.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusPermanentRedirect, "/login")
    })

	// Necessary for limiting
	router.ForwardedByClientIP = true

	// Enable session storage
	store := cookie.NewStore([]byte(configuration.SessionKey))
	
	if os.Getenv("RELEASE") == "production" {
		store.Options(sessions.Options{
			Domain:   configuration.Host,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		})

		gin.SetMode(gin.ReleaseMode)
	} else {
		store.Options(sessions.Options{
			Domain:   configuration.Host,
			Path:     "/",
		})
	}

	// Configure session storage
	router.Use(sessions.Sessions("session", store))

	// Load our HTML templates
	router.HTMLRender = loadTemplates("templates")

	// Load our static assets.
	router.Static("/assets", "dist")

	//
	//  Route Initialisation.
	//
	login.Init(router)
	browse.Init(router)
	play.Init(router)

	// Setup our 404
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	// If we're running in lambda, proxy. Otherwise, host on port 5000.
	if _, ok := os.LookupEnv("AWS_EXECUTION_ENV"); ok {
		// Connect to the ginadapter.
		ginLambda = ginadapter.New(router)
	} else {
		// Spin up on port :5000.
		router.Use(limit.MaxAllowed(3))

		err := router.Run(":5000")

		if err != nil {
			panic(err)
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepathx.Glob(templatesDir + "/views/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and views/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(strings.ReplaceAll(include, "templates/views/", ""), files...)
	}
	return r
}

//
//  This function is here to proxy any API Gateway request objects to Gin.
//
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

// Main Function.
func main() {
	lambda.Start(Handler)
}
