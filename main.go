package main

import (
	"fmt"
	"log"
	"net/url"

	"flag"
)

func main() {
	flag.Parse()
	config := loadConfig()
	overrideConfigFromFlags(&config)

	if cachedData, err := loadCachedTokenData(); err == nil && !*skipCache {
		fmt.Println(cachedData.AccessToken)
		return
	}

	state := createStateParam(16)

	authReqURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		config.AuthURL,
		url.QueryEscape(config.ClientID),
		url.QueryEscape(config.RedirectURI),
		url.QueryEscape(config.Scope),
		url.QueryEscape(state),
	)

	fmt.Println("Please open the following URL in your browser to continue: " + authReqURL)

	codeChan := make(chan string)
	go startCallbackServer(codeChan, state)

	authCode := <-codeChan

	data, err := exchangeCodeForToken(config, authCode)
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}

	if err := cacheTokenData(data); err != nil {
		log.Fatalf("Failed to cache credentials: %v", err)
	}

	fmt.Println(data.AccessToken)
}
