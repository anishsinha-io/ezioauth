package main

import (
	"fmt"
	"net/url"
	"os"

	"flag"
	"github.com/fatih/color"
	"golang.design/x/clipboard"
)

var (
	// config is the application configuration
	config appConfig

	mLog = color.RGB(119, 221, 119).Add(color.Bold)
	wLog = color.RGB(233, 236, 107).Add(color.Bold)
	eLog = color.RGB(244, 54, 76).Add(color.Bold)
)

func main() {
	flag.Parse()

	config = loadAppConfig()
	overrideConfigFromFlags(&config)

	if cachedData, err := loadCachedTokenData(); err == nil && !config.SkipCache {
		refreshed, err := refreshToken(cachedData.RefreshToken)
		if err != nil {
			wLog.Printf("Failed to refresh token (%v). Clearing cache...\n", err)
			os.Remove("credentials.json")
		} else {
			cacheTokenData(refreshed)

			clipboard.Write(clipboard.FmtText, []byte(refreshed.AccessToken))
			mLog.Println("Token copied to clipboard")
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

	mLog.Println("Please open the following URL in your browser to continue:\n" + authReqURL)

	codeChan := make(chan string)
	go startCallbackServer(codeChan, state)

	authCode := <-codeChan

	data, err := exchangeCodeForToken(authCode)
	if err != nil {
		eLog.Printf("Failed to exchange code for token: %v", err)
		os.Exit(1)
	}

	if err := cacheTokenData(data); err != nil {
		eLog.Printf("Failed to cache credentials: %v", err)
		os.Exit(1)
	}

	clipboard.Write(clipboard.FmtText, []byte(data.AccessToken))
	mLog.Println("Token copied to clipboard")
}
