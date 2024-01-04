package tasai

import (
	"os"

	"github.com/labstack/echo/v4"
)

type EchoApp struct {
	*echo.Echo
}

func NewEchoApp() *EchoApp {
	return &EchoApp{echo.New()}
}

func (e *EchoApp) Start(localPort ...string) error {
	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	// it's local
	if env != "dev" && env != "prod" {
		if localPort == nil || len(localPort) == 0 {
			localPort = []string{":8080"}
		}
		return e.Echo.Start(localPort[0])
	}
	return e.Echo.Start(port)
}
