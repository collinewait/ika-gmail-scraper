package oauth

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://ika-gmail-scraper-backend.herokuapp.com/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	client := getClient(googleOauthConfig, code)
	gmailService := getGmailService(client)
	log.Println(gmailService)

	io.WriteString(w, "Can call the api now, ") // nolint
}

func getClient(config *oauth2.Config, code string) *http.Client {
	token, err := googleOauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return config.Client(context.Background(), token)
}

func getGmailService(client *http.Client) *gmail.Service {
	service, err := gmail.New(client) // nolint
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return service
}