/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/gertd/go-pluralize"
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
		sql := generateSql(tableName, fields)
		fmt.Println(sql)
		generateMigrationFile(tableName, sql)
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
	Name         string
	Type         string
	IsUnique     bool
	IsRequired   bool
	IsPrimaryKey bool
	Default      string
	IsID         bool
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
			fieldType = strings.ReplaceAll(fieldType, "!", "")
		}

		if strings.Contains(fieldType, "^") {
			field.IsUnique = true
			fieldType = strings.ReplaceAll(fieldType, "^", "")
		}

		field.Type = fieldType

		newFields = append(newFields, field)
	}

	return newFields
}

func generateSqlField(field Field) string {

	var sqlString strings.Builder

	sqlString.WriteString(field.Name + " ")
	sqlString.WriteString(sqlType[field.Type])

	// if field.IsID {
	// 	sqlString.WriteString(fmt.Sprintf(sqlType[field.Type], pluralize.NewClient().Singular(tableName)))
	// } else {
	// 	sqlString.WriteString(sqlType[field.Type])
	// }

	if field.IsRequired {
		sqlString.WriteString(" NOT NULL")
	}

	if field.IsUnique {
		sqlString.WriteString(" UNIQUE")
	}

	if field.IsPrimaryKey {
		sqlString.WriteString(" PRIMARY KEY")
	}

	if field.Default != "" {
		sqlString.WriteString(" DEFAULT " + field.Default)
	}

	sqlString.WriteString(";")
	return sqlString.String()
}

func generateSql(tableName string, fields []Field) string {
	tmpl, err := template.New("main").Parse(sqlTemplate)

	if err != nil {
		fmt.Println(err)
		cobra.CheckErr("Failed to parse template")
	}

	sqlFields := make([]string, 0)

	hasID := false
	hasCreatedAt := false
	// hasUpdatedAt := false

	for _, field := range fields {
		switch field.Name {
		case "id":
			hasID = true
		case "created_at":
			hasCreatedAt = true
		case "updated_at":
			// hasUpdatedAt = true
		}
	}

	if !hasID {
		fieldString := generateSqlField(Field{
			Name:         "id",
			Type:         "nanoid",
			IsPrimaryKey: true,
			Default:      fmt.Sprintf("nanoid('%s_', 25)", pluralize.NewClient().Singular(tableName)),
		})
		sqlFields = append(sqlFields, fieldString)
	}

	for _, field := range fields {
		fieldString := generateSqlField(field)
		sqlFields = append(sqlFields, fieldString)
	}

	if !hasCreatedAt {
		sqlFields = append(sqlFields, generateSqlField(Field{
			Name:    "created_at",
			Type:    "datetime",
			Default: "NOW()",
		}))
	}

	// if !hasUpdatedAt {
	// 	sqlFields = append(sqlFields, generateSqlField(Field{
	// 		Name:    "updated_at",
	// 		Type:    "datetime",
	// 		Default: "NOW()",
	// 	}))
	//
	// }

	data := struct {
		TableName string
		Fields    []string
	}{
		TableName: tableName,
		Fields:    sqlFields,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)

	return buf.String()

}

var sqlType = map[string]string{
	"string":    "VARCHAR(255)",
	"text":      "TEXT",
	"int":       "INTEGER",
	"int32":     "INTEGER",
	"int64":     "BIGINT",
	"float":     "REAL",
	"float32":   "REAL",
	"float64":   "DOUBLE PRECISION",
	"bool":      "BOOLEAN",
	"boolean":   "BOOLEAN",
	"uuid":      "UUID",
	"date":      "DATE",
	"datetime":  "TIMESTAMP",
	"time":      "TIME",
	"json":      "JSON",
	"jsonb":     "JSONB",
	"bytea":     "BYTEA",
	"serial":    "SERIAL",
	"bigserial": "BIGSERIAL",
	"nanoid":    "TEXT",
}

func generateMigrationFile(tableName string, sqlContent string) {
	migratonDir := os.Getenv("MIGRATION_DIR")
	if migratonDir == "" {
		cobra.CheckErr("No migration directory provided.")
	}

	_, err := os.ReadDir(migratonDir)

	if err != nil {
		cobra.CheckErr("The migration directory wasn't found.")
	}

	timestamp := time.Now().Format("20060102150405")

	upName := fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, tableName)
	downName := fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, tableName)

	rootPath := migratonDir + "/"
	if err := createFile(rootPath+upName, sqlContent); err != nil {
		cobra.CheckErr(err)
	}

	if err := createFile(rootPath+upName, sqlContent); err != nil {
		cobra.CheckErr(err)
	}
	if err := createFile(rootPath+downName, ""); err != nil {
		cobra.CheckErr(err)
	}
}

func createFile(path, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
