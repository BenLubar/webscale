package db_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/BenLubar/webscale/db"
)

func TestTxRollback(t *testing.T) {
	db.InitForTesting()

	tx0, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin", err)
	defer tx0.Cancel()

	dropTableIfExists := db.Prepare(`drop table if exists testing_rollback;`)

	_, err = tx0.Exec(dropTableIfExists)
	checkFatalError(t, "drop table", err)

	createTable, err := tx0.Prepare(`create table testing_rollback (
	id bigserial,
	name text
);`)
	checkFatalError(t, "prepare create table", err)

	_, err = tx0.Exec(createTable)
	checkFatalError(t, "create table", err)

	err = tx0.Commit()
	checkFatalError(t, "commit", err)

	tx1, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin", err)
	defer tx1.Cancel()

	insert := db.Prepare(`insert into testing_rollback (name) values ($1) returning id;`)

	var helloID int64
	err = tx1.QueryRow(insert, "Hello").Scan(&helloID)
	checkFatalError(t, "insert", err)

	err = tx1.Commit()
	checkFatalError(t, "commit", err)

	tx2, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin 2", err)
	defer tx2.Cancel()

	get := db.Prepare(`select name from testing_rollback where id = $1;`)

	var hello string
	err = tx2.QueryRow(get, helloID).Scan(&hello)
	checkFatalError(t, "select", err)

	if hello != "Hello" {
		t.Errorf("expected %q, got %q", "Hello", hello)
	}

	var worldID int64
	err = tx2.QueryRow(insert, "World").Scan(&worldID)
	checkFatalError(t, "insert", err)

	err = tx2.Rollback()
	checkFatalError(t, "rollback", err)

	tx3, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin 3", err)
	defer tx3.Cancel()

	var world string
	err = tx3.QueryRow(get, worldID).Scan(&world)
	if err != sql.ErrNoRows {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = tx3.Exec(dropTableIfExists)
	checkFatalError(t, "drop table", err)

	err = tx3.Commit()
	checkFatalError(t, "commit", err)
}

func TestFailPrepare(t *testing.T) {
	db.InitForTesting()

	failed := make(chan struct{})

	oldPrepareFatal := *db.TestingPrepareFatal
	*db.TestingPrepareFatal = func(message string, done chan<- struct{}, wg *sync.WaitGroup) {
		close(failed)
		close(done)
		wg.Done()
	}
	defer func() {
		*db.TestingPrepareFatal = oldPrepareFatal
	}()

	succeeded := make(chan struct{}, 2)

	stmt := db.Prepare(`this is a syntax error`)
	go func() {
		stmt.Wait()
		select {
		case <-failed:
		default:
			t.Errorf("Wait finished unexpectedly")
			succeeded <- struct{}{}
		}
	}()
	go func() {
		db.WaitAll()
		select {
		case <-failed:
		default:
			t.Errorf("WaitAll finished unexpectedly")
			succeeded <- struct{}{}
		}
	}()

	select {
	case <-failed:
	case <-succeeded:
	}

	func() {
		const sentinel = "HELLO, WORLD"
		defer func() {
			if r := recover(); r != sentinel {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		oldPrepareFatal(sentinel, nil, nil)
	}()
}

func TestInvalidQuery(t *testing.T) {
	db.InitForTesting()

	tx, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin", err)
	defer tx.Cancel()

	stmt, err := tx.Prepare(`select now() > $1;`)
	checkFatalError(t, "prepare", err)

	rows, err := tx.Query(stmt, "Hello")
	if err == nil {
		t.Errorf("unexpected success")
	}
	if rows != nil {
		rows.Close()
	}
}
