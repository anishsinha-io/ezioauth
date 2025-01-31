package main

import (
	"context"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

// config represents the application configuration. It is a global variable initialized by `cmd`
var config appConfig

// cmd represents the CLI command that initializes and validates the application configuration
var cmd *cli.Command = &cli.Command{
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config-file",
			Usage:   "Path to the config file",
			Sources: cli.EnvVars("CONFIG_FILE"),
		},

		&cli.StringFlag{
			Name:        "server-auth-url",
			Usage:       "URL of the OpenID Connect authentication server",
			Sources:     cli.EnvVars("SERVER_AUTH_URL"),
			Destination: &config.Server.AuthURL,
		},

		&cli.StringFlag{
			Name:        "server-token-url",
			Usage:       "Token URL of the OpenID Connect authentication server",
			Sources:     cli.EnvVars("SERVER_TOKEN_URL"),
			Destination: &config.Server.TokenURL,
		},

		&cli.StringFlag{
			Name:        "server-client-id",
			Usage:       "ID of the client to authenticate against",
			Sources:     cli.EnvVars("SERVER_CLIENT_ID"),
			Destination: &config.Server.ClientID,
		},

		&cli.StringFlag{
			Name:        "server-client-secret",
			Usage:       "Secret of the client to authenticate against",
			Sources:     cli.EnvVars("SERVER_CLIENT_SECRET"),
			Destination: &config.Server.ClientSecret,
		},

		&cli.StringFlag{
			Name:        "server-redirect-uri",
			Usage:       "URL to redirect to after authentication",
			Sources:     cli.EnvVars("SERVER_REDIRECT_URI"),
			Destination: &config.Server.RedirectURI,
		},

		&cli.StringFlag{
			Name:        "server-scope",
			Usage:       "OpenID scope",
			Sources:     cli.EnvVars("SERVER_SCOPE"),
			Destination: &config.Server.Scope,
		},

		&cli.StringFlag{
			Name:        "callback-server-port",
			Usage:       "Port to listen on for the callback server",
			Sources:     cli.EnvVars("CALLBACK_SERVER_PORT"),
			Destination: &config.CallbackServerPort,
		},

		&cli.StringFlag{
			Name:        "credentials-cache",
			Usage:       "Path to where the app should cache the credentials",
			Sources:     cli.EnvVars("CREDENTIALS_CACHE"),
			Local:       true,
			Destination: &config.CredentialsCache,
		},

		&cli.BoolFlag{
			Name:        "skip-cache",
			Value:       false,
			Usage:       "Skip the cache and force a new token exchange",
			Sources:     cli.EnvVars("SKIP_CACHE"),
			Destination: &config.SkipCache,
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
				path, err := xdg.CacheFile("oauth-cli/credentials.json")
				if err != nil {
					return ctx, err
				}
				config.CredentialsCache = path
			}
		}

		if config.SkipCache == false {
			config.SkipCache = c.SkipCache
		}

		return ctx, nil
	},

	Action: func(ctx context.Context, command *cli.Command) error {
		return config.validate()
	},
}
