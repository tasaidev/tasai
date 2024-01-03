//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

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
	connection := os.Getenv("CONNECTION_STRING")
	// test using locally defined postgres url
	pg, err := tasai.Postgres(connection)
	if err != nil {
		panic(err)
	}
	defer pg.Close()
	db := bun.NewDB(pg, pgdialect.New())
	ctx := context.Background()
	res, err := db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	user := &User{Name: "admin"}
	res, err = db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		panic(err)
	}
}
