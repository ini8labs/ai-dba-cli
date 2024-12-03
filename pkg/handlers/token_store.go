package handlers

import (
	"encoding/json"
	"fmt"
	"os"
)

const SecretsFileName = "secrets.json"

type TokenStore struct {
	filePath string
}

// NewTokenStore creates a new FileTokenStore
func NewFileTokenStore(filePath string) *TokenStore {
	return &TokenStore{
		filePath: filePath,
	}
}

// SaveToken saves a token to a file
func (f *TokenStore) SaveToken(key, token string) error {
	// Ensure the directory exists
	dir := f.filePath
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Define the full file path
	// TODO: deprecate
	filePath := dir + "/" + SecretsFileName // This will be the file to store token data
	// filePath := dir // This will be the file to store token data

	// Read the existing data from the file, if it exists
	data := map[string]string{}
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(file, &data); err != nil {
			return err
		}
	}

	// Add the new token
	data[key] = token

	// Write the data back to the file
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Save the data to the file inside the tokens directory
	return os.WriteFile(filePath, file, 0644)
}

// GetToken retrieves a token from the file
func (f *TokenStore) GetToken(key string) (string, error) {
	data := map[string]string{}
	file, err := os.ReadFile(f.filePath + "/" + SecretsFileName)
	if err != nil {
		return "", fmt.Errorf("Readfile error %v", err)
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return "", err
	}

	token, ok := data[key]
	if !ok {
		return "", os.ErrNotExist
	}
	return token, nil
}
