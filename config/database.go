package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		"root:motherjoseph1412*@tcp(localhost:3306)/user_api?parseTime=true",
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
