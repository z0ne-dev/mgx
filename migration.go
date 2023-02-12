package mgx

import (
	"context"
	"fmt"
)

type MigrationFunc func(context.Context, Commands) error

type Migration interface {
	fmt.Stringer
	Run(context.Context, Commands) error
}

type migrationFuncWrapper struct {
	name string
	fn   MigrationFunc
}

func (m *migrationFuncWrapper) Run(ctx context.Context, tx Commands) error {
	return m.fn(ctx, tx)
}

func (m *migrationFuncWrapper) String() string {
	return m.name
}

func NewMigration(name string, fn MigrationFunc) Migration {
	return &migrationFuncWrapper{
		name: name,
		fn:   fn,
	}
}

func NewRawMigration(name, sql string) Migration {
	return &migrationFuncWrapper{
		name: name,
		fn:   func(ctx context.Context, tx Commands) error { _, err := tx.Exec(ctx, sql); return err },
	}
}
