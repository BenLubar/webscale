package internal_test

import (
	"testing"

	"github.com/BenLubar/webscale/internal"
	"github.com/pkg/errors"
)

func TestImpossibleError(t *testing.T) {
	internal.ImpossibleError(nil) // should not panic

	expectErr := errors.New("[expected error]")
	defer func() {
		if r := recover(); r != expectErr {
			t.Errorf("unexpected panic: %#v", r)
		}
	}()

	internal.ImpossibleError(expectErr) // should panic
}
