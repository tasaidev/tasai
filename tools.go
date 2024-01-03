//go:build tools
// +build tools

package tools

import (
	_ "github.com/Khan/genqlient/generate"
	_ "github.com/uptrace/bun"
	_ "github.com/uptrace/bun/dialect/pgdialect"
)
