// migrator.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the Apache 2.0 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: Apache-2.0

package mgx

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

const defaultTableName = "__migrations"

// Migrator is the migrator implementation
type Migrator struct {
	tableName  string
	logger     Logger
	migrations []Migration
}

// Option sets options such migrations or table name.
type Option func(*Migrator)

// TableName creates an option to allow overriding the default table name
func TableName(tableName string) Option {
	return func(m *Migrator) {
		m.tableName = tableName
	}
}

// Logger interface
type Logger interface {
	Log(msg string, data map[string]any)
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(msg string, data map[string]any)

// Log implements Logger interface
func (f LoggerFunc) Log(msg string, data map[string]any) {
	f(msg, data)
}

func defaultLogger(msg string, data map[string]any) {
	log.Println(msg, data)
}

// Log creates an option to allow overriding the stdout logging
func Log(logger Logger) Option {
	return func(m *Migrator) {
		m.logger = logger
	}
}

// Migrations creates an option with provided migrations
func Migrations(migrations ...Migration) Option {
	return func(m *Migrator) {
		m.migrations = migrations
	}
}

// New creates a new migrator instance
func New(opts ...Option) (*Migrator, error) {
	m := &Migrator{
		logger:    LoggerFunc(defaultLogger),
		tableName: defaultTableName,
	}
	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

// Migrate applies all available migrations
func (m *Migrator) Migrate(ctx context.Context, db Conn) error {
	// count applied migrations
	count, err := m.countApplied(ctx, db, m.tableName)
	if err != nil {
		return err
	}

	if count > len(m.migrations) {
		return ErrTooManyAppliedMigrations
	}

	m.logger.Log("Running missing migrations...", map[string]any{"missing": len(m.migrations) - count})

	// plan migrations
	for idx, migration := range m.migrations[count:] {
		err2 := m.applyMigration(ctx, db, migration, idx+count)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func (m *Migrator) applyMigration(ctx context.Context, db Conn, migration Migration, version int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("migrator: error while starting transaction: %w", err)
	}

	if errMigration := m.runMigration(ctx, tx, migration, version); errMigration != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			panic(fmt.Errorf("migrator: error while rolling back transaction: %w", errRollback))
		}
		return fmt.Errorf("migrator: error while running migrations: %w", errMigration)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("migrator: failed to commit transaction: %w", err)
	}
	return nil
}

// Pending returns all pending (not yet applied) migrations
func (m *Migrator) Pending(ctx context.Context, db Conn) ([]Migration, error) {
	count, err := m.countApplied(ctx, db, m.tableName)
	if err != nil {
		return nil, err
	}
	return m.migrations[count:len(m.migrations)], nil
}

func (m *Migrator) countApplied(ctx context.Context, db Conn, tableName string) (int, error) {
	// create migrations table if doesn't exist
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT8 NOT NULL,
			version VARCHAR(255) NOT NULL,
			PRIMARY KEY (id)
		);
	`, m.tableName)); err != nil {
		return 0, err
	}

	// count applied migrations
	var count int
	row := db.QueryRow(ctx, fmt.Sprintf("SELECT count(*) FROM %s", tableName))

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Migrator) runMigration(ctx context.Context, db pgx.Tx, migration Migration, id int) (err error) {
	m.logger.Log("applying migration", map[string]any{
		"id":   id,
		"name": migration.String(),
	})

	start := time.Now()
	if err = migration.Run(ctx, db); err != nil {
		return fmt.Errorf("error executing golang migration %s: %w", migration.String(), err)
	}

	insertVersion := fmt.Sprintf("INSERT INTO %s (id, version) VALUES ($1, $2)", m.tableName)
	if _, err = db.Exec(ctx, insertVersion, id, migration.String()); err != nil {
		return fmt.Errorf("error updating migration versions: %w", err)
	}
	duration := time.Since(start)

	m.logger.Log("applied migration", map[string]any{
		"id":   id,
		"name": migration.String(),
		"took": duration,
	})

	return err
}
