package tasai

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tasaidev/tasai/generated"
)

func Postgres(localUri ...string) (*sql.DB, error) {
	serviceId := os.Getenv("TASAI_SERVICE_ID")
	environmentId := os.Getenv("TASAI_ENVIRONMENT_ID")
	token := os.Getenv("TASAI_TOKEN")
	if serviceId != "" && environmentId != "" && token != "" {
		client := login(token)
		res, err := generated.CreatePostgresNeonResource(context.TODO(), client, serviceId, environmentId)
		if err != nil {
			return nil, err
		}
		if res != nil && res.CreatePostgresNeonResource.Id != "" {
			//success
			// weird I know but it's the only way to do this and prevent the serviec from spinning up
			os.Exit(42)
		}
		return nil, fmt.Errorf("failed to create postgres resource")
	}
	env := os.Getenv("ENV")
	// it's local don't use connections to deployed resources
	if env != "dev" && env != "prod" {
		if localUri == nil || len(localUri) == 0 {
			localUri = []string{"postgres://root:password@localhost:5432/postgres?sslmode=disable"}
		}
		return sql.Open("pgx", localUri[0])
	}
	// it's deployed use the connection string
	connectionString := os.Getenv("TASAI_POSTGRES_CONNECTION_STRING")
	if connectionString == "" {
		return nil, fmt.Errorf("failed to get postgres connection string")
	}
	return sql.Open("pgx", connectionString)
}
