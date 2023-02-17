// migration.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the Apache 2.0 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: Apache-2.0

// Package mgx is a simple migration tool for pgx
package mgx

import (
	"context"
	"fmt"
)

// MigrationFunc is a wrapper around a function so that it implements the Migration interface.
type MigrationFunc func(context.Context, Commands) error

// Migration is the migration interface
type Migration interface {
	fmt.Stringer
	Run(context.Context, Commands) error
}

type migrationFuncWrapper struct {
	fn   MigrationFunc
	name string
}

// Run the migration
func (m *migrationFuncWrapper) Run(ctx context.Context, tx Commands) error {
	return m.fn(ctx, tx)
}

// String returns the name of the migration
func (m *migrationFuncWrapper) String() string {
	return m.name
}

// NewMigration creates a migration from a function.
func NewMigration(name string, fn MigrationFunc) Migration {
	return &migrationFuncWrapper{
		name: name,
		fn:   fn,
	}
}

// NewRawMigration creates a migration from a raw SQL string.
func NewRawMigration(name, sql string) Migration {
	return &migrationFuncWrapper{
		name: name,
		fn:   func(ctx context.Context, tx Commands) error { _, err := tx.Exec(ctx, sql); return err },
	}
}
