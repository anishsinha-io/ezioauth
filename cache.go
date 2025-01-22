package main

import (
	"encoding/json"
	"os"
)

// cacheTokenData saves the token data to `credentials.json`.
func cacheTokenData(data tokenData) error {
	file, err := os.OpenFile("credentials.json", os.O_CREATE|os.O_WRONLY, 0600)
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

// loadCachedTokenData loads the token data from `credentials.json`.
func loadCachedTokenData() (tokenData, error) {
	file, err := os.Open("credentials.json")
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
