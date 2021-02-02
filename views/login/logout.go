package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLogout(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("user-id")
	session.Delete("state")

	err := session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/login")
}


