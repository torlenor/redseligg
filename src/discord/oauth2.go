package discord

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const htmlIndex = `<html><body>
<a href="/DiscordLogin">Log in with Discord</a>
</body></html>
`

type oauth2Handler struct {
	discordOauthConfig oauth2.Config
}

func (o *oauth2Handler) handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func (o *oauth2Handler) handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	url := o.discordOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *oauth2Handler) handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Errorf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	response, err := o.discordOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Errorf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Debugf("Joined a new server via OAuth2 flow. Response: %s", response)

	fmt.Fprint(w, "Successfully added the bot to the server! You can close the window now!")
}

func (o *oauth2Handler) startOAuth2Handler() {
	log.Infoln("Starting DiscordBot OAuth2 handler on http://localhost:8080/DiscordLogin")
	http.HandleFunc("/", o.handleMain)
	http.HandleFunc("/DiscordLogin", o.handleDiscordLogin)
	http.HandleFunc("/cb", o.handleDiscordCallback)

	log.Debug(http.ListenAndServe(":8080", nil))
}

func createOAuth2Handler(config oauth2.Config) *oauth2Handler {
	handler := &oauth2Handler{config}
	return handler
}
