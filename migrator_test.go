// migrator_test.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the Apache 2.0 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: Apache-2.0

package mgx_test

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/z0ne-dev/mgx/v2"
)

var migrations = []mgx.Migration{
	mgx.NewRawMigration("raw migration: single query, create table foo", "CREATE TABLE foo (id INT PRIMARY KEY)"),
	mgx.NewRawMigration(
		"raw migration: multi query, create table bar, alter table bar",
		"CREATE TABLE bar (id INT PRIMARY KEY); ALTER TABLE bar ADD COLUMN name VARCHAR(255)",
	),
	mgx.NewMigration("fn migration", func(ctx context.Context, cmd mgx.Commands) error {
		if _, err := cmd.Exec(ctx, "CREATE TABLE foobar (id INT PRIMARY KEY)"); err != nil {
			return err
		}
		if _, err := cmd.Exec(ctx, "INSERT INTO foobar (id) VALUES (1)"); err != nil {
			return err
		}
		return nil
	}),
}

var _ mgx.Logger = (*TestLogger)(nil)

type TestLogger struct {
	logged bool
}

func (t *TestLogger) Log(_ string, _ map[string]any) {
	t.logged = true
}

func connectToDatabase(t *testing.T) *pgx.Conn {
	t.Helper()

	// create db connection
	url := os.Getenv("POSTGRES")
	if url == "" {
		t.Fatal("POSTGRES env variable is not set")
	}

	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(context.Background(), "DROP SCHEMA public CASCADE"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(context.Background(), "CREATE SCHEMA public"); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMigrate(t *testing.T) {
	db := connectToDatabase(t)
	defer func(db *pgx.Conn, ctx context.Context) {
		_ = db.Close(ctx)
	}(db, context.Background())

	// create migrator
	migrator, err := mgx.New(mgx.Migrations(migrations...))
	if err != nil {
		t.Fatal(err)
	}

	// run migrator
	err = migrator.Migrate(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMigrateWithCustomLogger(t *testing.T) {
	db := connectToDatabase(t)
	defer func(db *pgx.Conn, ctx context.Context) {
		_ = db.Close(ctx)
	}(db, context.Background())

	l := new(TestLogger)

	// create migrator
	migrator, err := mgx.New(mgx.Log(l))
	if err != nil {
		t.Fatal(err)
	}

	// run migrator
	err = migrator.Migrate(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}

	if !l.logged {
		t.Fatal("custom logger was not called")
	}
}

func TestMigrateWithCustomTableName(t *testing.T) {
	db := connectToDatabase(t)
	defer func(db *pgx.Conn, ctx context.Context) {
		_ = db.Close(ctx)
	}(db, context.Background())

	// create migrator
	tableName := "custom_table_name"
	migrator, err := mgx.New(mgx.TableName(tableName))
	if err != nil {
		t.Fatal(err)
	}

	// run migrator
	err = migrator.Migrate(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}

	// check if table exists
	var rows int
	if err := db.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM "+tableName,
	).Scan(&rows); err != nil {
		t.Fatal(err)
	}
}

func TestPending(t *testing.T) {
	db := connectToDatabase(t)
	defer func(db *pgx.Conn, ctx context.Context) {
		_ = db.Close(ctx)
	}(db, context.Background())

	// create migrator
	migrator, err := mgx.New(mgx.Migrations(migrations...))
	if err != nil {
		t.Fatal(err)
	}

	pending, err := migrator.Pending(context.Background(), db)
	if err == nil && len(pending) != len(migrations) {
		t.Fatalf("there should be %d pending migrations, only %d found", len(migrations), len(pending))
	}

	if err != nil {
		t.Fatal(err)
	}
}
