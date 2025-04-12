package optional

import (
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestOptional(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		if err := json.Unmarshal([]byte(`null`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		testutil.RequireEqual(t, false, foo.hasValue)
		testutil.RequireEqual(t, "", foo.value)

		testutil.RequireNoError(t, foo.Validate())
		testutil.RequireEqual(t, true, foo.validated)
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`-24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		testutil.RequireEqual(t, true, foo.hasValue)
		testutil.RequireEqual(t, -24, foo.value) // won't be validated during unmarshalling

		testutil.RequireError(t, foo.Validate())
		testutil.RequireEqual(t, false, foo.validated)
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		testutil.RequireEqual(t, true, foo.hasValue)
		testutil.RequireEqual(t, 24, foo.value)

		testutil.RequireNoError(t, foo.Validate())
		testutil.RequireEqual(t, true, foo.validated)
	})
}
