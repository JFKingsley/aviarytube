package login

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jfkingsley/aviarytube/configuration"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func OAuthRedirect(c *gin.Context) {
	// If logged in, redirect

	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")

	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		session.AddFlash("Invalid session state")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		session.AddFlash("Login failed. Please try again")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		session.AddFlash("Error unmarshalling OAuth response")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	if !configuration.AuthorizedUsers[u.Email] {
		session.AddFlash("This user is not whitelisted.")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/browse")
}