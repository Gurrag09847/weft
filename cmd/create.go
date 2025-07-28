/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	templateRepo := "https://github.com/Gurrag09847/weft/tree/main/templates"

	repo, err := git.PlainClone(projectName, false, &git.CloneOptions{
		URL:      templateRepo,
		Progress: os.Stdout,
		Depth:    1,
	})

	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	worktree, err := repo.Worktree()

		return fmt.Errorf("failed to get worktree: %w", err)
	}

	sparseCheckoutPath := fmt.Sprintf("templates/%s", templateName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:                      plumbing.ZeroHash,
		Sparse                     CheckoutDirectories: []string{sparseCheckoutPath},
	})

	if err != nil {
		return fmt.Errorf("failed to checkout sparse directory: %w", err)
	}

	// Move the template files to the root of the project directory
		templateDir := filepath.Join(projectName, "templates", templateName)
	err := moveContents(templateDir, projectName); err != nil {
	turn fmt.Errorf("failed to move template contents: %w", err)
		
	
	// Clean up - remove templates directory and .git
	// Clean up - remove templates directory and .git
	err := os.RemoveAll(templatesDir); err != nil {
	turn fmt.Errorf("failed to remove templates directory: %w", err)
		
	
	gitDir := filepath.Join(projectName, ".git")
	err := os.RemoveAll(gitDir); err != nil {
	turn fmt.Errorf("failed to remove .git directory: %w", err)
		
	
	return nil
	

func moveContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if err := os.Rename(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}
