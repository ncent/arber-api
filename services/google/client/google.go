package clients

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	var GOOGLE_OAUTH_CLIENT_ID, _ = os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
	var GOOGLE_OAUTH_CLIENT_SECRET, _ = os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET")
	var GOOGLE_OAUTH_ENDPOINT_TOKEN_URL, _ = os.LookupEnv("GOOGLE_OAUTH_ENDPOINT_TOKEN_URL")

	GoogleClient = NewGoogleService()
	GoogleOAuthConfig = &GoogleConfig{
		GOOGLE_OAUTH_CLIENT_ID,
		GOOGLE_OAUTH_CLIENT_SECRET,
		GOOGLE_OAUTH_ENDPOINT_TOKEN_URL,
		nil,
		&oauth2.Config{
			GOOGLE_OAUTH_CLIENT_ID,
			GOOGLE_OAUTH_CLIENT_SECRET,
			oauth2.Endpoint{
				"",
				GOOGLE_OAUTH_ENDPOINT_TOKEN_URL,
				0,
			},
			"",
			nil,
		},
		&clientcredentials.Config{
			GOOGLE_OAUTH_CLIENT_ID,
			GOOGLE_OAUTH_CLIENT_SECRET,
			GOOGLE_OAUTH_ENDPOINT_TOKEN_URL,
			nil,
			nil,
			0,
		},
	}
}

var GoogleClient *GoogleService
var GoogleOAuthConfig *GoogleConfig
