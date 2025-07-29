/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var fieldRegex = regexp.MustCompile("^([a-zA-Z_][a-zA-Z0-9_]*):([a-zA-Z][a-zA-Z0-9]*(?:\\\\([0-9,\\\\s]*\\\\))?)(!)?([\\\\^])?$")

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// modelName := args[0]
		tableName := args[0]
		fields := extractFields(args[1:])
		fmt.Println("schema called", tableName, fields)
		generateSql(tableName, fields)
	},
}

func init() {

	genCmd.AddCommand(schemaCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schemaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schemaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Field struct {
	Name       string
	Type       string
	IsUnique   bool
	IsRequired bool
}

func extractFields(fields []string) []Field {
	newFields := make([]Field, 0)
	for _, f := range fields {

		matches := fieldRegex.FindStringSubmatch(f)

		if matches == nil {
			cobra.CheckErr("Field format must be name:type[!][^]")
		}

		parts := strings.Split(f, ":")

		if len(parts) != 2 {
			cobra.CheckErr("Field format must be name:type[!][^]")
		}

		fieldType := parts[1]
		field := Field{
			Name: parts[0],
		}
		if strings.Contains(fieldType, "!") {
			field.IsRequired = true
			fieldType = fieldType[:len(fieldType)-1]
		}

		if strings.Contains(fieldType, "^") {
			field.IsUnique = true

			fieldType = fieldType[:len(fieldType)-1]
		}

		field.Type = fieldType

		newFields = append(newFields, field)
	}

	return newFields
}

func generateSqlField(field Field) string {

	var sqlString strings.Builder

	sqlString.WriteString(field.Name)
	sqlString.WriteString(field.Type)

	return sqlString.String()
}

func generateSql(tableName string, fields []Field) {
	tmpl, err := template.New("main").Parse(sqlTemplate)

	if err != nil {
		fmt.Println(err)
		cobra.CheckErr("Failed to parse template")
	}

	sqlFields := make([]string, 0)

	for _, field := range fields {
		fieldString := generateSqlField(field)
		sqlFields = append(sqlFields, fieldString)
	}

	data := struct {
		TableName string
		Fields    []string
	}{
		TableName: tableName,
		Fields:    sqlFields,
	}

	err = tmpl.Execute(os.Stdout, data)

}

var sqlType map[string]string = map[string]string{
	"string": "VARCHAR(255)",
}
