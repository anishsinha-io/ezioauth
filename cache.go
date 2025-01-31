package main

import (
	"encoding/json"
	"os"
)

// cacheTokenData saves the token data to the cache file path configured
//
// The token data is saved as JSON that looks like:
//
//	{
//	  "access_token": "<a signed jwt>",
//	  "refresh_token": "<a signed jwt>"
//	}
//
// When new data is saved, the old data is overwritten. The path where the credentials are
// cached can be configured by the `CREDENTIALS_CACHE` environment variable or the
// `--credentials-cache` flag or the "credentials_cache" field in the configuration file.
func cacheTokenData(data tokenData) error {
	file, err := os.OpenFile(config.CredentialsCache, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

// loadCachedTokenData loads the token data from the configured cache file path
func loadCachedTokenData() (tokenData, error) {
	file, err := os.Open(config.CredentialsCache)
	if err != nil {
		return tokenData{}, err
	}
	defer file.Close()

	var data tokenData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return tokenData{}, err
	}

	return data, nil
}
