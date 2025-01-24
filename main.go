package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"flag"
)

var (
	// config is the application configuration
	config appConfig
)

func main() {
	flag.Parse()

	config = loadAppConfig()
	overrideConfigFromFlags(&config)

	if cachedData, err := loadCachedTokenData(); err == nil && !config.SkipCache {
		refreshed, err := refreshToken(cachedData.RefreshToken)
		if err != nil {
			fmt.Printf("Failed to refresh token: %v\n", err)
			os.Remove("credentials.json")
		} else {
			cacheTokenData(refreshed)
			fmt.Println(refreshed.AccessToken)
			return
		}
	}

	state := createStateParam(16)

	authReqURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		config.Server.AuthURL,
		url.QueryEscape(config.Server.ClientID),
		url.QueryEscape(config.Server.RedirectURI),
		url.QueryEscape(config.Server.Scope),
		url.QueryEscape(state),
	)

	fmt.Println("Please open the following URL in your browser to continue: " + authReqURL)

	codeChan := make(chan string)
	go startCallbackServer(codeChan, state)

	authCode := <-codeChan

	data, err := exchangeCodeForToken(authCode)
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}

	if err := cacheTokenData(data); err != nil {
		log.Fatalf("Failed to cache credentials: %v", err)
	}

	fmt.Println(data.AccessToken)
}
