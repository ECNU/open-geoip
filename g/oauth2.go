package g

import (
	"fmt"
	"golang.org/x/oauth2"
)

var OauthConfig *oauth2.Config

func InitOauth() {
	fmt.Printf("2222%v\n\n", Config().Oauth.ClientId)
	OauthConfig = &oauth2.Config{
		ClientID:     Config().Oauth.ClientId,
		ClientSecret: Config().Oauth.ClientSecret,
		RedirectURL:  Config().Oauth.RedirectURL,
		Scopes:       Config().Oauth.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Config().Oauth.SsoAddr,
			TokenURL: Config().Oauth.TokenAddr,
		},
	}
}

//
//http.HandleFunc("/", handleMain)
//http.HandleFunc("/login", handleLogin)
//http.HandleFunc("/callback", handleCallback)
//
//log.Println("Server started on http://localhost:8080")
//log.Fatal(http.ListenAndServe(":8080", nil))
