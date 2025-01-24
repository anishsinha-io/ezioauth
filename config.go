package main

import (
	"encoding/json"
	"flag"
	"os"
)

var (
	// configFile is the path to the configuration file.
	configFile = flag.String("config", "", "Path to the config file")

	// authURL, tokenURL, clientID, clientSecret, redirectURI, and scope are the OpenID Connect configuration parameters.
	authURL          = flag.String("auth-url", "", "OpenID authentication URL")
	tokenURL         = flag.String("token-url", "", "OpenID token URL")
	clientID         = flag.String("client-id", "", "OpenID client ID")
	clientSecret     = flag.String("client-secret", "", "OpenID client secret")
	redirectURI      = flag.String("redirect-uri", "", "OpenID redirect URI")
	scope            = flag.String("scope", "", "OpenID scopes")
	credentialsCache = flag.String("credentials-cache", "", "Path to where the app should cache the credentials")
	skipCache        = flag.Bool("skip-cache", false, "Skip the cache and force a new token exchange")
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

// appConfig represents the application configuration.
type appConfig struct {
	Server           serverConfig `json:"server"`
	CredentialsCache string       `json:"credentials_cache"`
	SkipCache        bool         `json:"skip_cache"`
}

// tokenData represents the returned data as a result of a successful OAuth 2.0 authorization code flow
type tokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// loadAppConfig loads the application configuration from the configuration file.
func loadAppConfig() appConfig {
	c := appConfig{}
	if *configFile != "" {
		file, err := os.Open(*configFile)
		if err != nil {
			eLog.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&c); err != nil {
			eLog.Println(err)
			os.Exit(1)
		}
	}

	return c
}

// overrideConfigFromFlags overrides the configuration parameters from the command line flags.
func overrideConfigFromFlags(c *appConfig) {
	if *authURL != "" {
		c.Server.AuthURL = *authURL
	}
	if *tokenURL != "" {
		c.Server.TokenURL = *tokenURL
	}
	if *clientID != "" {
		c.Server.ClientID = *clientID
	}
	if *clientSecret != "" {
		c.Server.ClientSecret = *clientSecret
	}
	if *redirectURI != "" {
		c.Server.RedirectURI = *redirectURI
	}
	if *scope != "" {
		c.Server.Scope = *scope
	}
	if *credentialsCache != "" {
		c.CredentialsCache = *credentialsCache
	}
	if *skipCache {
		c.SkipCache = *skipCache
	}
}
