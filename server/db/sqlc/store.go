package db

import (
	"database/sql"
	"fmt"
	"peer-talk/config"

	_ "github.com/lib/pq"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	nativeDb *sql.DB
}

func NewSQLStore() (*SQLStore, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Env.DbUser,
		config.Env.DbPassword,
		config.Env.DbHost,
		config.Env.DbPort,
		config.Env.DbName,
		config.Env.DbSslMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &SQLStore{
		nativeDb: db,
		Queries:  New(db),
	}, nil
}
