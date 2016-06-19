// Package schema handles database structure management for #webscale.
package schema // import "github.com/BenLubar/webscale/db/internal/schema"

import (
	"database/sql"
	"strconv"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// Upgrade handles the initialization of the database.
func Upgrade(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin schema upgrade transaction")
	}
	defer tx.Rollback()

	// Grab a transaction-level advisory lock. The number isn't important
	// as long as it stays the same for all #webscale instances.
	if _, err = tx.Exec(`select pg_advisory_xact_lock(1);`); err != nil {
		return errors.Wrap(err, "acquire upgrade lock")
	}

	if _, err = tx.Exec(`create table if not exists schema_changes (
	version bigint primary key,
	description text not null,
	applied_at timestamp with time zone not null default now()
);`); err != nil {
		return errors.Wrap(err, "create schema_changes table")
	}

	// version is the last index of all that was applied to the database.
	var version sql.NullInt64

	err = tx.QueryRow(`select max(sc.version) from schema_changes as sc;`).Scan(&version)
	if err != nil {
		return errors.Wrap(err, "get schema version")
	}

	if !version.Valid {
		version.Int64 = -1
	}

	const target = int64(len(all)) - 1
	if version.Int64 > target {
		return errors.Errorf("unknown schema version %d - is this database from a newer version of #webscale?", version.Int64)
	}

	if version.Int64 == target {
		return nil
	}

	for i, change := range all[version.Int64+1:] {
		id := int64(i) + version.Int64 + 1
		if _, err = tx.Exec(change.query); err != nil {
			if pe, ok := err.(*pq.Error); ok && pe.Position != "" {
				if pos, ne := strconv.Atoi(pe.Position); ne == nil && pos > 0 {
					near := change.query[pos-1:]
					if len(near) > 100 {
						near = near[:100]
					}
					err = errors.Wrapf(err, "near %q", near)
				}
			}
			return errors.Wrapf(err, "apply change #%d: %s", id, change.description)
		}
		if _, err = tx.Exec(`insert into schema_changes (version, description) values ($1, $2);`, id, change.description); err != nil {
			return errors.Wrapf(err, "record change #%d: %s", id, change.description)
		}
	}

	return errors.Wrap(tx.Commit(), "commit schema changes")
}
