package db_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/BenLubar/webscale/db"
)

func checkFatalError(t *testing.T, message string, err error) {
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

func TestPostgresConstraint(t *testing.T) {
	db.InitForTesting()

	t.Run("Empty", func(t *testing.T) {
		testPostgresConstraint(t, testPostgresConstraintEmpty)
	})

	t.Run("Hello1", func(t *testing.T) {
		testPostgresConstraint(t, testPostgresConstraintHello1)
	})

	t.Run("Hello2", func(t *testing.T) {
		testPostgresConstraint(t, testPostgresConstraintHello2)
	})
}

func testPostgresConstraint(t *testing.T, f func(*testing.T, *db.Tx, *db.Stmt)) {
	tx, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin", err)
	defer tx.Cancel()

	createTable, err := tx.Prepare(`create temporary table test_constraint (
	id bigserial,
	name text not null
		constraint name_not_empty check (name <> '')
		constraint name_unique unique
);`)
	checkFatalError(t, "prepare create table", err)
	_, err = tx.Exec(createTable)
	checkFatalError(t, "create temporary table", err)

	insert, err := tx.Prepare(`insert into test_constraint (name) values ($1) returning id;`)
	checkFatalError(t, "prepare insert", err)

	f(t, tx, insert)
}

func testPostgresConstraintEmpty(t *testing.T, tx *db.Tx, insert *db.Stmt) {
	var id0 int64
	err := tx.QueryRow(insert, "").Scan(&id0)
	if err == nil {
		t.Errorf("expected error for empty name")
	} else if !db.IsConstraint(err, "name_not_empty") {
		t.Errorf("incorrect error for empty name: %v", err)
	}

	if db.IsConstraint(err, "name_unique") {
		t.Errorf("db.IsConstraint returned true for the wrong constraint")
	}
}

func testPostgresConstraintHello1(t *testing.T, tx *db.Tx, insert *db.Stmt) {
	var id1 int64
	err := tx.QueryRow(insert, "Hello").Scan(&id1)
	if err != nil {
		t.Errorf("insert failed unexpectedly: %v", err)
	} else if id1 <= 0 {
		t.Errorf("insert succeeded but ID is %d", id1)
	}

	if db.IsConstraint(err, "") {
		t.Errorf("db.IsConstraint returned true for %v", err)
	}
}

func testPostgresConstraintHello2(t *testing.T, tx *db.Tx, insert *db.Stmt) {
	testPostgresConstraintHello1(t, tx, insert)

	var id2 int64
	err := tx.QueryRow(insert, "Hello").Scan(&id2)
	if err == nil {
		t.Errorf("expected error for duplicate name")
	} else if !db.IsConstraint(err, "name_unique") {
		t.Errorf("incorrect error for duplicate name: %v", err)
	} else if err.(interface {
		Timeout() bool
	}).Timeout() {
		t.Errorf("unexpected timeout error: %v", err)
	} else if err.(interface {
		Temporary() bool
	}).Temporary() {
		t.Errorf("unexpected temporary error: %v", err)
	}
}

func TestPostgresErrorDetails(t *testing.T) {
	db.InitForTesting()

	tx, err := db.Begin(context.Background())
	checkFatalError(t, "db.Begin", err)
	defer tx.Cancel()

	createTable, err := tx.Prepare(`create temporary table test_errors (
	id bigserial,
	birthday date
);`)
	checkFatalError(t, "prepare create table", err)

	_, err = tx.Exec(createTable)
	checkFatalError(t, "create temporary table", err)

	insert, err := tx.Prepare(`insert into test_errors (birthday) values ($1) returning id;`)
	checkFatalError(t, "prepare insert", err)

	goBirthday := time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC)
	var id0 int64
	err = tx.QueryRow(insert, goBirthday).Scan(&id0)
	checkFatalError(t, "insert Go", err)

	googleBirthday := time.Date(1998, time.September, 27, 0, 0, 0, 0, time.UTC)
	var id1 int64
	err = tx.QueryRow(insert, googleBirthday).Scan(&id1)
	checkFatalError(t, "insert Google", err)

	get, err := tx.Prepare(`select id, birthday from test_errors order by birthday desc;`)
	checkFatalError(t, "prepare select", err)

	func() {
		rows, err := tx.Query(get)
		checkFatalError(t, "select", err)
		defer rows.Close()

		var id int64
		var birthday time.Time

		if !rows.Next() {
			t.Fatalf("rows.Next returned false")
		}
		err = rows.Scan(&id, &birthday)
		checkFatalError(t, "rows.Scan", err)
		if id != id0 || !birthday.Equal(goBirthday) {
			t.Errorf("(%d, %v) != (%d, %v)", id, birthday, id0, goBirthday)
		}

		if !rows.Next() {
			t.Fatalf("rows.Next returned false")
		}
		err = rows.Scan(&id, &birthday)
		checkFatalError(t, "rows.Scan", err)
		if id != id1 || !birthday.Equal(googleBirthday) {
			t.Errorf("(%d, %v) != (%d, %v)", id, birthday, id1, googleBirthday)
		}

		if rows.Next() {
			t.Fatalf("rows.Next returned true")
		}

		checkFatalError(t, "rows.Err", rows.Err())
	}()

	var id2 int64
	err = tx.QueryRow(insert, "Cake!").Scan(&id2)
	if err == nil {
		t.Fatalf("unexpected nil error")
	}

	s := err.Error()
	if !strings.HasPrefix(s, "pq: invalid input syntax for type date: \"Cake!\"\n\nStack trace:\n") {
		t.Errorf("error does not start with expected string\n%s", s)
	}
	if !strings.Contains(s, "\n\nQuery:\n\ninsert into test_errors (birthday) values ($1) returning id;\n\nPrepared at:\n") {
		t.Errorf("error does not contain expected query\n%s", s)
	}
	if !strings.HasSuffix(s, "\n\nDriver details\nError code: 22007 (invalid_datetime_format)\nMessage: invalid input syntax for type date: \"Cake!\"") {
		t.Errorf("error does not end with expected string\n%s", s)
	}
}
