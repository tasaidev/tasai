//go:build ignore

package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tasaidev/tasai"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID   int64 `bun:",pk,autoincrement"`
	Name string
}

func main() {
	app := tasai.NewEchoApp(
		tasai.WithName("test"),
		tasai.WithProject("test"),
		tasai.WithLocalPort(8080),
		tasai.WithDevInstance(tasai.Instance{
			Maximum: 1,
		}),
		tasai.WithProdInstance(tasai.Instance{
			Minimum: 1,
			Maximum: 3,
		}),
		tasai.WithPostgres("postgres://root:password@localhost:5433/postgres?sslmode=disable"),
	)
	// auto closes db connection
	pg, err := app.Postgres()
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(pg, pgdialect.New())
	ctx := context.Background()
	_, err = db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}
	user := &User{Name: "admin"}
	_, err = db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		panic(err)
	}

	app.GET("/", func(c echo.Context) error {
		return c.String(200, "hello world")
	})

	app.GET("/users", func(c echo.Context) error {
		var users []User
		err := db.NewSelect().Model(&users).Scan(ctx)
		if err != nil {
			return err
		}
		return c.JSON(200, users)
	})

	app.Start()
}
