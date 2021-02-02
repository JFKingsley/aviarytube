package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/configuration"
	"github.com/oklog/ulid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"math/rand"
	"net/http"
	"time"
)

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     configuration.OAuthClientID,
		ClientSecret: configuration.OAuthSecret,
		RedirectURL:  configuration.Protocol + "://" + configuration.Host + "/login/oauth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func GetLogin(c *gin.Context) {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	state := id.String()

	session := sessions.Default(c)
	session.Set("state", state)

	flash := ""
	if flashes := session.Flashes(); len(flashes) > 0 {
		flash = flashes[0].(string)
	}

	err := session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}

	c.HTML(http.StatusOK, "login/index.html", gin.H{
		"Link": conf.AuthCodeURL(state),
		"Flash": flash,
	})
}

