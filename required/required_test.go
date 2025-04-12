package required

import (
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestRequired(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		if err := json.Unmarshal([]byte(`null`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		testutil.RequireEqual(t, false, foo.hasValue)
		testutil.RequireEqual(t, "", foo.value)

		testutil.RequireError(t, foo.Validate())
		testutil.RequireEqual(t, false, foo.validated)

		testutil.RequirePanic(t, func() { foo.Get() })
		testutil.RequirePanic(t, func() { foo.MarshalJSON() })
		testutil.RequirePanic(t, func() { foo.MarshalText() })
		testutil.RequirePanic(t, func() { foo.Value() })
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

		testutil.RequirePanic(t, func() { foo.Get() })
		testutil.RequirePanic(t, func() { foo.MarshalJSON() })
		testutil.RequirePanic(t, func() { foo.MarshalText() })
		testutil.RequirePanic(t, func() { foo.Value() })
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

		testutil.RequireNoPanic(t, func() { foo.Get() })
		testutil.RequireNoPanic(t, func() { foo.MarshalJSON() })
		testutil.RequireNoPanic(t, func() { foo.MarshalText() })
		testutil.RequireNoPanic(t, func() { foo.Value() })
	})
}
