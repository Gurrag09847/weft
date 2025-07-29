/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Creating project....")

		if args[0] == "" {
			fmt.Println("provide a proejct name")
		}
		createProject(args[0], "gin_postgres_htmx")

		fmt.Println("Created project!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createProject(projectName, templateName string) error {
	// Create temporary directory for full clone
	tempDir := fmt.Sprintf("%s-temp", projectName)
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Clone the full repository to temp directory
	templateRepo := "https://github.com/Gurrag09847/weft.git"
	_, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      templateRepo,
		Progress: os.Stdout,
		Depth:    1, // Shallow clone for efficiency
	})
	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	// Create the project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Copy only the specific template
	templateDir := filepath.Join(tempDir, "templates", templateName)
	if err := copyDir(templateDir, projectName); err != nil {
		return fmt.Errorf("failed to copy template contents: %w", err)
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
