package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	BaseURL    = "https://dba-api-xxxvi-guk4.onxplorx.app"
	UIURL      = "https://dba-fe-xxix-tsos.onxplorx.app"
	WebhookURL = BaseURL + "/v1/data"
	Binary     = "dba.exe"
	// Binary     = "dba-linux-amd64"
	// Binary = "dba-darwin-arm64"

	// BaseURL    = "http://localhost:3000"
	// WebhookURL = BaseURL + "/v1/data"
	// Binary     = "dba.exe"
)

var rootCmd = &cobra.Command{
	Use:     Binary,
	Short:   "Data Base Analyser",
	Long:    fmt.Sprintf("A database analysis tool that supports postgresql database.\nAnalyze databases running locally or you can also check out the web app at %s for more features and insights.\nSupported connection string formats:\n\tPostgreSQL: postgresql://user:pass@host:port/dbname", UIURL),
	Example: "\t" + Binary + " login -e <email> -p <password>\n\t" + Binary + " analyse -c <postgres_connection_string>",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Welcome to dba!")

		cmd.Help()
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// to add custom flags
func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// rootCmd.SetHelpCommand(&cobra.Command{
	// 	Hidden: true, // Hides the help command
	// })

	// rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Example:\n\n\t%s analyse -c postgresql://user:pass@localhost:5432/dbname\n\t%s analyse --connection-string postgresql://user:pass@127.0.0.1:5432/dbname\n\nNote: `localhost` and `120.0.0.1` can be used interchangeably.", Binary, Binary)
	// })

	rootCmd.Flags().BoolP("help", "h", false, "Help about any command")
}
