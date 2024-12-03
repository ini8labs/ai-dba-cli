package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dba.exe",
	Short: "Data Base Analyser",
	Long: `A database analysis tool that supports postgresql database.
Supported connection string formats:
  PostgreSQL: postgresql://user:pass@host:5432/dbname`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infoln("Welcome to dba!")
	},
}

func Execute() {
	helpFlag := false

	// Check for -help or -h flag
	for _, arg := range os.Args[1:] {
		if arg == "-help" || arg == "-h" {
			helpFlag = true
			break
		}
	}

	if helpFlag {
		logrus.Info("Help will be provided shortly.")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// to add custom flags
func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help message")
}
