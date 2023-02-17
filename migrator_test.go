package mgx_test

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/z0ne-dev/mgx"
	"os"
	"testing"
)

var migrations = []mgx.Migration{
	mgx.NewRawMigration("raw migration: single query, create table foo", "CREATE TABLE foo (id INT PRIMARY KEY)"),
	mgx.NewRawMigration("raw migration: multi query, create table bar, alter table bar", "CREATE TABLE bar (id INT PRIMARY KEY); ALTER TABLE bar ADD COLUMN name VARCHAR(255)"),
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

func TestMigrate(t *testing.T) {
	// create db connection
	url := os.Getenv("POSTGRES")
	if url == "" {
		t.Fatal(errors.New("POSTGRES env variable is not set"))
	}

	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: replace with your db connection

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
