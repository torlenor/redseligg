package discord

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const htmlIndex = `<html><body>
<a href="/DiscordLogin">Log in with Discord</a>
</body></html>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	url := discordOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := discordOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Token: %s\n", token)
}

func doAuthOAuth2() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/DiscordLogin", handleDiscordLogin)
	http.HandleFunc("/cb", handleDiscordCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
