/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weft",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weft.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getMigrationDir() string {
	migratonDir := os.Getenv("MIGRATION_DIR")
	if migratonDir == "" {
		cobra.CheckErr("No migration directory provided.")
	}

	_, err := os.ReadDir(migratonDir)

	if err != nil {
		cobra.CheckErr("The migration directory wasn't found.")
	}

	return migratonDir
}

func getDatabaseURL() string {
	url := os.Getenv("PSQL_DSN")
	if url == "" {
		cobra.CheckErr("No database url provided.")
	}
	return url
}

func loadENV() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func createFile(path, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
