package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// startCallbackServer starts a simple HTTP server that listens for the OAuth 2.0 authorization code.
func startCallbackServer(codeChan chan<- string, expectedState string) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")

		if code == "" || state != expectedState {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Authorization successful. You can close this tab.")

		codeChan <- code
	})

	log.Fatal(http.ListenAndServe(":8666", nil))
}

// exchangeCodeForToken exchanges the authorization code for an access token.
func exchangeCodeForToken(authCode string) (tokenData, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", config.Server.RedirectURI)
	data.Set("client_id", config.Server.ClientID)
	data.Set("client_secret", config.Server.ClientSecret)

	req, err := http.NewRequest("POST", config.Server.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return tokenData{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return tokenData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return tokenData{}, fmt.Errorf("token exchange failed: %s", body)
	}

	var respData tokenData
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return tokenData{}, err
	}

	return respData, nil
}

// refreshToken refreshes the access token using the given refresh token
func refreshToken(refreshToken string) (tokenData, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", config.Server.ClientID)
	data.Set("client_secret", config.Server.ClientSecret)

	req, err := http.NewRequest("POST", config.Server.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return tokenData{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return tokenData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return tokenData{}, fmt.Errorf("token refresh failed: %s", resp.Status)
	}

	var respData tokenData

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return tokenData{}, err
	}

	return respData, nil
}

// createStateParam generates a random string of the specified length. It can be used as the OAuth 2.0 state parameter.
func createStateParam(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
