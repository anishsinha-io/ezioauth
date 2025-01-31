package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/adrg/xdg"
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

// validate validates the server configuration.
// TODO: Improve validation error messages
func (s *serverConfig) validate() error {
	if s.AuthURL == "" || s.TokenURL == "" ||
		s.ClientID == "" || s.ClientSecret == "" ||
		s.RedirectURI == "" || s.Scope == "" {
		return fmt.Errorf("server configuration is missing required fields")
	}
	return nil
}

// appConfig represents the application configuration.
type appConfig struct {
	Server             serverConfig `json:"server"`
	CredentialsCache   string       `json:"credentials_cache"`
	SkipCache          bool         `json:"skip_cache"`
	CallbackServerPort string       `json:"callback_server_port"`
}

// validate validates the application configuration.
func (a *appConfig) validate() error {
	if err := a.Server.validate(); err != nil {
		return err
	}

	if a.CallbackServerPort == "" {
		return fmt.Errorf("callback server port is missing")
	}

	u, err := url.Parse(a.Server.RedirectURI)
	if err != nil {
		return fmt.Errorf("invalid redirect URI: %v", err)
	}

	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return fmt.Errorf("invalid redirect URI: %v", err)
	}

	if port != a.CallbackServerPort {
		return fmt.Errorf("redirect URI port does not match the callback server port")
	}

	return nil
}

// tokenData represents the returned data as a result of a successful OAuth 2.0 authorization code flow
type tokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// initXDGConfig initializes the configuration file in the XDG configuration directory.
func initXDGConfig() (string, error) {
	configFilePath, err := xdg.ConfigFile("oauth-cli/config.json")
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_RDWR, 0600)

	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	if fileInfo.Size() == 0 {
		_, err = file.WriteString("{}")
		if err != nil {
			return "", err
		}
		file.Seek(0, 0)
	}
	return configFilePath, nil
}

// initFromFile initializes the configuration from the given file.
func initFromFile(filePath string, config *appConfig) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return err
	}
	return nil
}
