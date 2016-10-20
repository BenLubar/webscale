package internal_test

import (
	"testing"

	"github.com/BenLubar/webscale/internal"
)

func recurseCallers(n int) internal.StackTrace {
	if n == 0 {
		return internal.Callers(0)
	}
	return recurseCallers(n - 1)
}

func TestCallersCount(t *testing.T) {
	n0 := len(recurseCallers(0))
	n4 := len(recurseCallers(4))
	n16 := len(recurseCallers(16))
	n64 := len(recurseCallers(64))

	if n4-n0 != 4 {
		t.Errorf("expected %d callers, got %d", 4, n4-n0)
	}
	if n16-n0 != 16 {
		t.Errorf("expected %d callers, got %d", 16, n16-n0)
	}
	if n64-n0 != 64 {
		t.Errorf("expected %d callers, got %d", 64, n64-n0)
	}
}

func TestInvalidProgramCounter(t *testing.T) {
	stack := internal.StackTrace{0}
	s := string(stack.AppendTo(nil))
	if s != "\n(0) (unknown function)+0x0" {
		t.Errorf("expected %q, got %q", "\n(0) (unknown function)+0x0", s)
	}
}
