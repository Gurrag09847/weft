/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Gurrag09847/weft/cmd/templates"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]
		generateDomain(domainName)
	},
}

func init() {
	genCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateDomain(name string) {
	handlerTmpl, err := template.New("handler").Parse(templates.HandlerTemplate)

	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to parse template: %w", err))
	}

	plural := pluralize.NewClient().Pluralize(name, 2, false)
	singular := pluralize.NewClient().Singular(name)

	modulePath := getModulePath()

	data := struct {
		PackageName  string
		DomainPlural string
		DomainName   string
		ModulePath   string
		DomainUpper  string
		Plural       string
	}{
		PackageName:  name,
		DomainPlural: cases.Title(language.English).String(plural),
		Plural:       cases.Title(language.English).String(name),
		DomainName:   name,
		ModulePath:   modulePath,
		DomainUpper:  cases.Title(language.English).String(singular),
	}

	var handlerBuf bytes.Buffer
	err = handlerTmpl.Execute(&handlerBuf, data)

	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to execute template: %w", err))
	}

	serviceTmpl, err := template.New("service").Parse(templates.ServiceTemplate)

	if err != nil {
		fmt.Errorf("failed to parse template: %w", err)
	}

	var serviceBuf bytes.Buffer
	err = serviceTmpl.Execute(&serviceBuf, data)

	if err != nil {
		fmt.Errorf("failed to execute template: %w", err)
	}

	_, err = os.ReadDir("./internal")

	if err != nil {
		cobra.CheckErr("No internal directory found.")
	}
	path := fmt.Sprintf("./internal/%s", name)

	err = os.MkdirAll(path, 0755)

	if err != nil {
		cobra.CheckErr(fmt.Errorf("Error creating directory: %w", err))
	}

	err = createFile(fmt.Sprintf("%s/handler.go", path), handlerBuf.String())
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Error creating files: %w", err))
	}

	err = createFile(fmt.Sprintf("%s/service.go", path), serviceBuf.String())
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Error creating files: %w", err))
	}

	fmt.Println("Created domain files!")
}

func getModulePath() string {
	file, err := os.ReadFile("go.mod")

	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to read go.mod: %w", err))
	}

	modulePath := strings.Split(strings.Split(string(file), "\n")[0], "module ")[1]

	return modulePath
}
