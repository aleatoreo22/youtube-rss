package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"youtube-rss.api.aleatoreo.com/package/database/sqlite"
)

type Database struct {
	Framework DatabaseFramework
}

type DatabaseFramework interface {
	Open() error
	Close() error
	DatabaseExists() bool
	Execute(string) (sql.Result, error)
	CreateTable(any) error
	Query(string, any) ([]any, error)
}

type QueryBuilder struct {
	strings.Builder
}

func NewQueryBuilder(query string) *QueryBuilder {
	qb := &QueryBuilder{}
	qb.WriteString(query)
	return qb
}

func (qb *QueryBuilder) AddParameter(parameter string, data any) {
	var dataString string

	switch v := data.(type) {

	case string:
		dataString = "'" + strings.ReplaceAll(v, "'", "''") + "'"

	case time.Time:
		dataString = "'" + v.UTC().Format(time.RFC3339) + "'"

	case *time.Time:
		if v == nil {
			dataString = "NULL"
		} else {
			dataString = "'" + v.UTC().Format(time.RFC3339) + "'"
		}

	case bool:
		if v {
			dataString = "1"
		} else {
			dataString = "0"
		}

	case nil:
		dataString = "NULL"

	default:
		dataString = fmt.Sprintf("%v", v)
	}

	query := strings.Replace(qb.String(), parameter, dataString, 1)

	qb.Reset()
	qb.WriteString(query)
}

func New(dbType string) (*Database, error) {
	var framework DatabaseFramework
	switch dbType {
	case "sqlite":
		framework = &sqlite.SQLite{}
	default:
		return nil, errors.New("unsupported database type")
	}
	database := &Database{Framework: framework}
	return database, nil
}

func (db *Database) CreateDatabase(models []any) error {
	if db.Framework.DatabaseExists() {
		return nil
	}
	err := (db.Framework).Open()
	if err != nil {
		return err
	}
	for _, model := range models {
		err = (db.Framework).CreateTable(model)
		if err != nil {
			return err
		}
	}
	(db.Framework).Close()
	return nil
}

func ConvertQueryResult[T any](result []any) ([]T, error) {
	var converted []T
	for _, v := range result {
		converted = append(converted, v.(T))
	}
	return converted, nil
}
