package tasai

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgres struct {
	db *sql.DB
}

func newPostgres(localUri string) (*postgres, error) {
	isdeployment := os.Getenv("TASAI_DEPLOYMENT")
	if isdeployment != "" {
		return nil, nil
	}
	env := os.Getenv("ENV")
	// it's local don't use connections to deployed resources
	if env != "dev" && env != "prod" {
		db, err := sql.Open("pgx", localUri)
		if err != nil {
			return nil, err
		}
		return &postgres{
			db,
		}, nil
	}
	// it's deployed use the connection string
	connectionString := os.Getenv("TASAI_POSTGRES_CONNECTION_STRING")
	if connectionString == "" {
		return nil, fmt.Errorf("failed to get postgres connection string")
	}
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	return &postgres{
		db,
	}, nil
}
