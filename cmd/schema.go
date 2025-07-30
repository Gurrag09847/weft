/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Gurrag09847/weft/cmd/templates"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

var fieldRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*):([a-zA-Z][a-zA-Z0-9_]*(?:\([0-9,\s]*\))?)([!^]{0,2})?(?:=(.*))?$`)

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
		loadENV()
		// modelName := args[0]
		tableName := args[0]
		fields := extractFields(args[1:])
		upSql, downSql := generateSql(tableName, fields)
		fmt.Println(upSql, downSql)
		generateMigrationFile(tableName, upSql, downSql)
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
	Name           string
	Type           string
	IsUnique       bool
	IsRequired     bool
	IsPrimaryKey   bool
	Default        string
	IsReference    bool
	ReferenceField string
}

func extractFields(fields []string) []Field {
	newFields := make([]Field, 0)
	for _, f := range fields {

		matches := fieldRegex.FindStringSubmatch(f)

		if matches == nil {
			cobra.CheckErr("Field format must be name:type[!][^][=default]")
		}

		fieldType := matches[2]
		field := Field{
			Name: matches[1],
		}
		modifiers := matches[3]
		if strings.Contains(modifiers, "!") {
			field.IsRequired = true
			fieldType = strings.ReplaceAll(fieldType, "!", "")
		}

		if strings.Contains(modifiers, "^") {
			field.IsUnique = true
			fieldType = strings.ReplaceAll(fieldType, "^", "")
		}

		field.Type = fieldType

		if matches[4] != "" && matches[2] != "reference" {
			field.Default = matches[4]
		} else if matches[2] == "reference" {
			field.IsReference = true
			field.ReferenceField = matches[4]
			field.Type = "text"
			field.Name = field.Name + "_id"
		}

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

	if field.Default != "" && !field.IsReference {
		sqlString.WriteString(" DEFAULT " + field.Default)
	}

	if field.IsReference {
		sqlString.WriteString(fmt.Sprintf(" REFERENCES %s(id)", field.ReferenceField))
	}

	return sqlString.String()
}

func generateSql(tableName string, fields []Field) (string, string) {
	upTmpl, err := template.New("main").Parse(templates.UpTemplate)

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
	err = upTmpl.Execute(&buf, data)

	downTempl, err := template.New("down").Parse(templates.DownTemplate)

	if err != nil {
		fmt.Println(err)
		cobra.CheckErr("Failed to parse template")
	}

	downData := struct {
		TableName string
	}{
		TableName: tableName,
	}

	var downBuf bytes.Buffer
	err = downTempl.Execute(&downBuf, downData)
	return buf.String(), downBuf.String()
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

func generateMigrationFile(tableName string, sqlUpContent, sqlDownContent string) {
	migratonDir := getMigrationDir()
	timestamp := time.Now().Format("20060102150405")

	upName := fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, tableName)
	downName := fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, tableName)

	rootPath := migratonDir + "/"
	if err := createFile(rootPath+upName, sqlUpContent); err != nil {
		cobra.CheckErr(err)
	}

	if err := createFile(rootPath+downName, sqlDownContent); err != nil {
		cobra.CheckErr(err)
	}
}
