package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	// serverConfigFile is the path to the server config file. An example is provided at `example-server.json`
	serverConfigFile = flag.String("server-config", "", "Path to the server config file")

	// authURL, tokenURL, clientID, clientSecret, redirectURI, and scope are the OpenID Connect configuration parameters.
	authURL      = flag.String("auth-url", "", "OpenID authentication URL")
	tokenURL     = flag.String("token-url", "", "OpenID token URL")
	clientID     = flag.String("client-id", "", "OpenID client ID")
	clientSecret = flag.String("client-secret", "", "OpenID client secret")
	redirectURI  = flag.String("redirect-uri", "", "OpenID redirect URI")
	scope        = flag.String("scope", "", "OpenID scopes")

	// skipCache is a flag that determines whether to skip loading cached credentials.
	skipCache = flag.Bool("skip-cache", false, "Skip loading cached credentials")
)

// serverConfig represents the OpenID Connect configuration parameters.
type serverConfig struct {
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

// tokenData represents the returned data as a result of a successful OAuth 2.0 authorization code flow
type tokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// loadConfig loads the server configuration from the server config file.
func loadConfig() serverConfig {
	c := serverConfig{}
	if *serverConfigFile != "" {
		file, err := os.Open(*serverConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&c); err != nil {
			log.Fatal(err)
		}
	}

	return c
}

// overrideConfigFromFlags overrides the server configuration parameters with the flags provided.
func overrideConfigFromFlags(config *serverConfig) {
	if *authURL != "" {
		config.AuthURL = *authURL
	}

	if *tokenURL != "" {
		config.TokenURL = *tokenURL
	}

	if *clientID != "" {
		config.ClientID = *clientID
	}

	if *clientSecret != "" {
		config.ClientSecret = *clientSecret
	}

	if *redirectURI != "" {
		config.RedirectURI = *redirectURI
	}

	if *scope != "" {
		config.Scope = *scope
	}
}
