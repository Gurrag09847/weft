/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			cobra.CheckErr("provide a migration name")
		}
		loadENV()
		generateMigrationFile(args[0])
	},
}

func init() {
	migrateCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateMigrationFile(name string) {
	migratonDir := getMigrationDir()
	timestamp := time.Now().Format("20060102150405")

	upName := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downName := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	rootPath := migratonDir + "/"
	if err := createFile(rootPath+upName, ""); err != nil {
		cobra.CheckErr(err)
	}

	if err := createFile(rootPath+downName, ""); err != nil {
		cobra.CheckErr(err)
	}
}
