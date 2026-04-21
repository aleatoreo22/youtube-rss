package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

const sqliteDatabaseFile = "file:database.db?_loc=auto&parseTime=true"

func (sqlite *SQLite) Open() error {
	db, err := sql.Open("sqlite3", sqliteDatabaseFile)
	if err != nil {
		return err
	}
	sqlite.db = db
	return nil
}

func (sqlite *SQLite) Close() error {
	return sqlite.db.Close()
}

func (db *SQLite) DatabaseExists() bool {
	_, err := os.Stat(sqliteDatabaseFile)
	return !errors.Is(err, os.ErrNotExist)
}

func GetStringLen(tag string) int {
	re := regexp.MustCompile(`strlen\s+(\d+)`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return 0
	}
	strlen, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}
	return strlen
}

func IsPrimaryKey(tag string) bool {
	return strings.Contains(tag, "key")
}

func IsUniqueColumn(tag string) bool {
	return strings.Contains(tag, "unique")
}

func relativeType(typeName string, len int) string {
	relativeType := ""
	switch typeName {
	case "string":
		if len > 0 {
			relativeType = "varchar"
		} else {
			relativeType = "TEXT"
		}
	case "int", "int32", "int64":
		relativeType = "INTEGER"
	case "float32", "float64":
		relativeType = "REAL"
	case "bool":
		relativeType = "BOOLEAN"
	case "Time":
		relativeType = "DATETIME"
	default:
		relativeType = "TEXT"
	}
	return relativeType
}

func (sqlite *SQLite) CreateTable(tableModel any) error {
	t := reflect.TypeOf(tableModel)
	tableName := t.Name()
	var columns []string
	var primaryKey []string
	var uniqueColumn []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("sql")
		var queryField strings.Builder
		queryField.WriteString(field.Name + " ")
		fieldLen := GetStringLen(tag)
		queryField.WriteString(relativeType(field.Type.Name(), fieldLen) + " ")
		if fieldLen > 0 {
			queryField.WriteString("(" + strconv.Itoa(fieldLen) + ")")
		}
		if IsPrimaryKey(tag) {
			queryField.WriteString(" NOT NULL ")
			primaryKey = append(primaryKey, field.Name)
		}
		if IsUniqueColumn(tag) {
			queryField.WriteString(" NOT NULL ")
			uniqueColumn = append(uniqueColumn, field.Name)
		}
		columns = append(columns, queryField.String())
	}
	sql := "CREATE TABLE IF NOT EXISTS " + tableName + " (" + strings.Join(columns, ",\n")
	if len(primaryKey) > 0 {
		sql += ", \n PRIMARY KEY (" + strings.Join(primaryKey, ", ") + ")"
	}
	if len(uniqueColumn) > 0 {
		sql += ", \n UNIQUE (" + strings.Join(uniqueColumn, ", ") + ")"
	}
	sql += ");"
	_, err := sqlite.Execute(sql)
	return err
}

func (sqlite *SQLite) Execute(query string) (sql.Result, error) {
	res, err := sqlite.db.Exec(query)
	return res, err
}

func (sqlite *SQLite) Query(query string, model any) ([]any, error) {
	rows, err := sqlite.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("db"); tag != "" {
			fieldMap[tag] = i
		} else {
			fieldMap[field.Name] = i
		}
	}

	var results []any

	for rows.Next() {
		element := reflect.New(t).Elem()

		pointers := make([]any, len(columns))
		tempTime := make(map[int]*string)

		for i, column := range columns {
			if fieldIndex, ok := fieldMap[column]; ok {
				field := element.Field(fieldIndex)
				if field.Type() == reflect.TypeOf(time.Time{}) {
					var temp string
					tempTime[i] = &temp
					pointers[i] = &temp
				} else {
					pointers[i] = field.Addr().Interface()
				}
			} else {
				var dummy any
				pointers[i] = &dummy
			}
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		for i, strPtr := range tempTime {
			if strPtr != nil && *strPtr != "" {
				fieldIndex := fieldMap[columns[i]]
				field := element.Field(fieldIndex)
				fmt.Println("RAW TIME:", *strPtr)

				parsed, err := parseSQLiteTime(*strPtr)
				if err == nil {
					field.Set(reflect.ValueOf(parsed))
				}
			}
		}

		results = append(results, element.Interface())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05 -0700 -0700", 
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s", value)
}