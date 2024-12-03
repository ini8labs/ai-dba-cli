package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CLI struct {
	BaseURL    string
	TokenStore TokenStore
}

// NewCLI initializes the CLI with the server's base URL
func NewCLI(baseURL string, tokenStore TokenStore) *CLI {
	return &CLI{
		BaseURL:    baseURL,
		TokenStore: tokenStore,
	}
}

// Login sends login credentials to the server and stores the JWT locally
func (cli *CLI) Login(email, password string) error {
	logrus.Debugln("Calling Login via CLI")
	url := fmt.Sprintf("%s/v1/users/login", cli.BaseURL)

	// Prepare the request payload
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", body)
	}

	// Parse the response
	var response struct {
		Message string `json:"message"`
		Token   string `json:"token"`
		User    struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("Login successful! Welcome %s. \n", response.User.Email)

	// Save the JWT token locally
	if err := cli.TokenStore.SaveToken("auth_token", response.Token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}
