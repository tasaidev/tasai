package tasai

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tasaidev/tasai/internal"
)

type EchoApp struct {
	*echo.Echo
	settings *appSettings
	postgres *internal.Postgres
	env      string
}

func NewEchoApp(opts ...appOption) *EchoApp {
	settings, err := processAndValidateOpts(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	isdeployment := os.Getenv("TASAI_DEPLOYMENT")
	if isdeployment != "" {
		// don't block by opening tcp listener
		b, err := json.Marshal(settings)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		os.Exit(0)
	}
	env := os.Getenv("ENV")
	if env != "dev" && env != "prod" {
		env = "local"
	}
	var postgres *internal.Postgres
	if settings.Postgres {
		pg, err := internal.NewPostgres(settings.Localuri)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		postgres = pg
	}
	return &EchoApp{
		echo.New(),
		settings,
		postgres,
		env,
	}
}

func (e *EchoApp) IsLocal() bool {
	return e.env == "local"
}

func (e *EchoApp) IsDev() bool {
	return e.env == "dev"
}

func (e *EchoApp) IsProd() bool {
	return e.env == "prod"
}

func (e *EchoApp) Postgres() (*sql.DB, error) {
	if e.postgres == nil || e.postgres.DB == nil {
		return nil, fmt.Errorf("postgres not configured")
	}
	return e.postgres.DB, nil
}

func (e *EchoApp) Start() error {
	if e.postgres != nil && e.postgres.DB != nil {
		defer e.postgres.DB.Close()
	}
	// it's local
	if e.IsLocal() {
		return e.Echo.Start(e.settings.Localport)
	}
	port := os.Getenv("PORT")
	return e.Echo.Start(port)
}
