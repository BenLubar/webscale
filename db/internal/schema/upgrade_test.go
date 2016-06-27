package schema_test

import (
	"testing"

	"github.com/BenLubar/webscale/db"
	"github.com/BenLubar/webscale/internal/testutils"
)

func TestUpgrade(t *testing.T) {
	if err := db.Init(*testutils.FlagDB); err != nil {
		t.Error(err)
	}
}
