package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
	"golang.design/x/clipboard"
)

// config represents the application configuration. It is a global variable initialized by `cmd`
var config appConfig

// author represents the author of the application
type author struct {
	name  string
	email string
}

func (a author) String() string {
	return a.name + " <" + a.email + ">"
}

// cmd represents the CLI command that initializes and validates the application configuration
var cmd *cli.Command = &cli.Command{
	Name:    "EzioAuth",
	Version: "0.1.0",
	Authors: []any{
		author{
			name:  "Anish",
			email: "anishsinha0128@gmail.com",
		},
	},
	Copyright:   "(c) 2025 EzioAuth authors",
	Usage:       "Retrieve an access token from an OpenID server using the authorization code flow",
	UsageText:   "ezioauth [options]",
	Description: "EzioAuth is a CLI utility that lets you retrieve an access token from an OpenID Connect server using the authorization code flow.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config-file",
			Usage:    "Path to the config file",
			Sources:  cli.EnvVars("CONFIG_FILE"),
			Category: "Server",
		},

		&cli.BoolFlag{
			Name:     "save-config",
			Value:    true,
			Usage:    "Save the configuration to the default config file",
			Sources:  cli.EnvVars("SAVE_CONFIG"),
			Category: "Server",
		},

		&cli.StringFlag{
			Name:        "server-auth-url",
			Usage:       "URL of the OpenID Connect authentication server",
			Sources:     cli.EnvVars("SERVER_AUTH_URL"),
			Destination: &config.Server.AuthURL,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "server-token-url",
			Usage:       "Token URL of the OpenID Connect authentication server",
			Sources:     cli.EnvVars("SERVER_TOKEN_URL"),
			Destination: &config.Server.TokenURL,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "server-client-id",
			Usage:       "ID of the client to authenticate against",
			Sources:     cli.EnvVars("SERVER_CLIENT_ID"),
			Destination: &config.Server.ClientID,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "server-client-secret",
			Usage:       "Secret of the client to authenticate against",
			Sources:     cli.EnvVars("SERVER_CLIENT_SECRET"),
			Destination: &config.Server.ClientSecret,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "server-redirect-uri",
			Usage:       "URL to redirect to after authentication",
			Sources:     cli.EnvVars("SERVER_REDIRECT_URI"),
			Destination: &config.Server.RedirectURI,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "server-scope",
			Usage:       "OpenID scope",
			Sources:     cli.EnvVars("SERVER_SCOPE"),
			Destination: &config.Server.Scope,
			Category:    "Server",
		},

		&cli.StringFlag{
			Name:        "callback-server-port",
			Usage:       "Port to listen on for the callback server",
			Sources:     cli.EnvVars("CALLBACK_SERVER_PORT"),
			Destination: &config.CallbackServerPort,
			Category:    "Global",
		},

		&cli.StringFlag{
			Name:        "credentials-cache",
			Usage:       "Path to where the app should cache the credentials",
			Sources:     cli.EnvVars("CREDENTIALS_CACHE"),
			Destination: &config.CredentialsCache,
			Category:    "Global",
		},

		&cli.BoolFlag{
			Name:     "skip-cache",
			Value:    false,
			Usage:    "Skip the cache and force a new token exchange",
			Sources:  cli.EnvVars("SKIP_CACHE"),
			Category: "Global",
		},
	},

	Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
		configFile := command.String("config-file")
		var c appConfig

		xdgConfigPath, err := initXDGConfig()
		if err != nil {
			return ctx, err
		}

		if configFile != "" {
			if err := initFromFile(configFile, &c); err != nil {
				return ctx, err
			}
		} else {

			if err := initFromFile(xdgConfigPath, &c); err != nil {
				return ctx, err
			}
		}

		if config.Server.AuthURL == "" {
			config.Server.AuthURL = c.Server.AuthURL
		}

		if config.Server.TokenURL == "" {
			config.Server.TokenURL = c.Server.TokenURL
		}

		if config.Server.ClientID == "" {
			config.Server.ClientID = c.Server.ClientID
		}

		if config.Server.ClientSecret == "" {
			config.Server.ClientSecret = c.Server.ClientSecret
		}

		if config.Server.RedirectURI == "" {
			config.Server.RedirectURI = c.Server.RedirectURI
		}

		if config.Server.Scope == "" {
			config.Server.Scope = c.Server.Scope
		}

		if config.CallbackServerPort == "" {
			config.CallbackServerPort = c.CallbackServerPort
		}

		if config.CredentialsCache == "" {
			if c.CredentialsCache != "" {
				config.CredentialsCache = c.CredentialsCache
			} else {
				path, err := xdg.CacheFile("ezioauth/credentials.json")
				if err != nil {
					return ctx, err
				}
				config.CredentialsCache = path
			}
		}

		return ctx, nil
	},

	Action: func(ctx context.Context, command *cli.Command) error {
		err := config.validate()
		if err != nil {
			return err
		}

		if command.Bool("save-config") {
			configFilePath, err := xdg.ConfigFile("ezioauth/config.json")
			if err != nil {
				return err
			}
			if err := config.save(configFilePath); err != nil {
				return err
			}
		}

		skipCache := command.Bool("skip-cache")

		return run(skipCache)
	},
}

// run runs the application
func run(skipCache bool) error {
	if cachedData, err := loadCachedTokenData(); err == nil && !skipCache {
		refreshed, err := refreshToken(cachedData.RefreshToken)
		if err != nil {
			fmt.Printf("Failed to refresh token (%v). Clearing cache...\n", err)
			os.Remove(config.CredentialsCache)
		} else {
			if err := cacheTokenData(refreshed); err != nil {
				return err
			}

			clipboard.Write(clipboard.FmtText, []byte(refreshed.AccessToken))
			fmt.Println("Token copied to clipboard")
			return nil
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
		return fmt.Errorf("Failed to exchange code for token: %v", err)
	}

	if err := cacheTokenData(data); err != nil {
		return fmt.Errorf("Failed to cache credentials: %v", err)
	}

	clipboard.Write(clipboard.FmtText, []byte(data.AccessToken))
	fmt.Println("Token copied to clipboard")
	return nil
}
