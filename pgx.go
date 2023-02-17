// pgx.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the Apache 2.0 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: Apache-2.0

package mgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Conn is the interface that wraps the basic connection methods.
type Conn interface {
	Commands
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Commands is the interface that wraps the basic sql command methods.
type Commands interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
