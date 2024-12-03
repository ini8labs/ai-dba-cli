package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ini8labs/ai-dba-cli/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to your account",
		Long:  `Login to your account to use the AI DBA CLI.`,
		RunE:  login,
	}
)

// for email and password, and adds it to the root command.
func init() {
	loginCmd.Flags().StringP("email", "e", "", "Email address")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	rootCmd.AddCommand(loginCmd)
}

func login(cmd *cobra.Command, args []string) error {

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return err
	} else if email == "-p" {
		return fmt.Errorf("email cannot be empty; please provide it using the --email or -e flag")
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return err
	}

	// Validate email and password inputs
	if email == "" {
		return fmt.Errorf("email cannot be empty; please provide it using the --email or -e flag")
	}
	if password == "" {
		return fmt.Errorf("password cannot be empty; please provide it using the --password or -p flag")
	}

	// TODO: implement login
	// Prepare the request payload
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := os.Getenv("BASE_URL")
	loginURL := fmt.Sprintf("%s/v1/users/login", url)

	// Send the POST request
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("login failed: %s. Check credentials and try again.", body)
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

	logrus.Infof("Login successful! Welcome %s.", response.User.Email)

	// Save the token in the config file
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	config.Token = response.Token
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	logrus.Debugln("Token saved successfully.")

	return nil
}
