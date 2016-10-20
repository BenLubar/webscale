package db_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/BenLubar/hook"
	"github.com/BenLubar/webscale/db"
	"github.com/pkg/errors"
)

func TestErrorTimeout(t *testing.T) {
	db.InitForTesting()

	ctx, cancel := context.WithTimeout(context.Background(), -time.Second)
	defer cancel()

	tx, err := db.Begin(ctx)
	if tx != nil {
		t.Errorf("tx is not nil")
		tx.Cancel()
	}

	if err == nil {
		t.Fatalf("err is nil")
	}

	cause := errors.Cause(err)
	if cause != context.DeadlineExceeded {
		t.Errorf("unexpected cause: %v", cause)
	}

	stack := err.(interface {
		StackTrace() errors.StackTrace
	}).StackTrace()
	if fn := fmt.Sprintf("%n", stack[0]); fn != "TestErrorTimeout" {
		t.Errorf("StackTrace says we're in %q", fn)
	}

	if !err.(interface {
		Timeout() bool
	}).Timeout() {
		t.Errorf("Timeout says no")
	}

	if !err.(interface {
		Temporary() bool
	}).Temporary() {
		t.Errorf("Temporary says no")
	}

	if strings.Contains(err.Error(), "\n\nDriver details\n") {
		t.Errorf("unexpected driver error: %v", err)
	}

	if db.IsConstraint(err, "timeout") {
		t.Errorf("unexpected constraint error: %v", err)
	}
}

func TestErrorFilterError(t *testing.T) {
	db.InitForTesting()

	actualApply := *db.TestingApplyFilterAppendErrorDetails
	defer func() {
		*db.TestingApplyFilterAppendErrorDetails = actualApply
	}()
	fakeRegister := hook.NewFilter(db.TestingApplyFilterAppendErrorDetails).(func(func([]byte, context.Context) ([]byte, error), int))

	tx, err := db.Begin(context.Background())
	if err != nil {
		t.Fatalf("db.Begin: %v", err)
	}
	defer tx.Cancel()

	_, err = tx.Prepare(`invalid syntax for a query`)
	if err == nil {
		t.Fatalf("no error from invalid query")
	}

	msg1 := err.Error()

	fakeRegister(func(buf []byte, ctx context.Context) ([]byte, error) {
		return buf, errors.New("[expected error]")
	}, 0)

	msg2 := err.Error()

	if msg1 == msg2 {
		t.Fatalf("error messages are equal: %v", msg1)
	}

	const sentinel = "\n\n!!! ERROR WHILE GETTING ERROR DETAILS !!!\n[expected error]"
	i := strings.Index(msg2, sentinel)
	if i == -1 {
		t.Fatalf("sentinel string not found: %v", msg2)
	}

	if msg2[:i]+msg2[i+len(sentinel):] != msg1 {
		t.Errorf("messages differ in an unexpected way:\n%q\n%q", msg1, msg2)
	}
}
