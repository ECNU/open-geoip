package g

import (
	"fmt"

	"golang.org/x/oauth2"
)

var Oauth2Config *oauth2.Config

func InitOauth2() {
	fmt.Printf("2222%v\n\n", Config().Oauth2.ClientId)
	Oauth2Config = &oauth2.Config{
		ClientID:     Config().Oauth2.ClientId,
		ClientSecret: Config().Oauth2.ClientSecret,
		RedirectURL:  Config().Oauth2.RedirectURL,
		Scopes:       Config().Oauth2.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Config().Oauth2.AuthAddr,
			TokenURL: Config().Oauth2.TokenAddr,
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
