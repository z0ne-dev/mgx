// error.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the Apache 2.0 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: Apache-2.0

package mgx

import "errors"

// ErrTooManyAppliedMigrations is returned when more migrations are applied than defined
var ErrTooManyAppliedMigrations = errors.New("too many applied migrations")
