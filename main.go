package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"golang.design/x/clipboard"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

	if cachedData, err := loadCachedTokenData(); err == nil && !config.SkipCache {
		refreshed, err := refreshToken(cachedData.RefreshToken)
		if err != nil {
			fmt.Printf("Failed to refresh token (%v). Clearing cache...\n", err)
			os.Remove(config.CredentialsCache)
		} else {
			cacheTokenData(refreshed)

			clipboard.Write(clipboard.FmtText, []byte(refreshed.AccessToken))
			fmt.Println("Token copied to clipboard")
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

	fmt.Println("Please open the following URL in your browser to continue [cmd+click]:\n" + authReqURL)

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

	clipboard.Write(clipboard.FmtText, []byte(data.AccessToken))
	fmt.Println("Token copied to clipboard")
}
