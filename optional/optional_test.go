package optional

import (
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestOptional(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		testutil.RequireNoError(t, json.Unmarshal([]byte(`null`), &foo))

		testutil.RequireEqual(t, false, foo.hasValue)
		testutil.RequireEqual(t, foo.hasValue, foo.HasValue())
		testutil.RequireEqual(t, "", foo.value)

		testutil.RequireNoError(t, foo.Validate())
		testutil.RequireEqual(t, true, foo.validated)

		testutil.RequirePanic(t, func() { foo.Must() })
		testutil.RequireNoPanic(t, func() { foo.Get() })
		testutil.RequireNoPanic(t, func() { foo.MarshalJSON() })
		testutil.RequireNoPanic(t, func() { foo.MarshalText() })
		testutil.RequireNoPanic(t, func() { foo.Value() })
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.RequireNoError(t, json.Unmarshal([]byte(`-24`), &foo))

		testutil.RequireEqual(t, true, foo.hasValue)
		testutil.RequireEqual(t, foo.hasValue, foo.HasValue())
		testutil.RequireEqual(t, -24, foo.value)

		testutil.RequireError(t, foo.Validate())
		testutil.RequireEqual(t, false, foo.validated)

		testutil.RequirePanic(t, func() { foo.Must() })
		testutil.RequirePanic(t, func() { foo.Get() })
		testutil.RequirePanic(t, func() { foo.MarshalJSON() })
		testutil.RequirePanic(t, func() { foo.MarshalText() })
		testutil.RequirePanic(t, func() { foo.Value() })
	})

	t.Run("nested invalid value", func(t *testing.T) {
		type Foo struct {
			Field Positive[int]
		}

		var foo Any[Foo]

		testutil.RequireNoError(t, json.Unmarshal([]byte(`{"field":-1}`), &foo))
		testutil.RequireError(t, foo.Validate())
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.RequireNoError(t, json.Unmarshal([]byte(`24`), &foo))

		testutil.RequireEqual(t, true, foo.hasValue)
		testutil.RequireEqual(t, foo.hasValue, foo.HasValue())
		testutil.RequireEqual(t, 24, foo.value)

		testutil.RequireNoError(t, foo.Validate())
		testutil.RequireEqual(t, true, foo.validated)

		testutil.RequireNoPanic(t, func() { foo.Must() })
		testutil.RequireNoPanic(t, func() { foo.Get() })
		testutil.RequireNoPanic(t, func() { foo.MarshalJSON() })
		testutil.RequireNoPanic(t, func() { foo.MarshalText() })
		testutil.RequireNoPanic(t, func() { foo.Value() })

		t.Run("reuse as invalid", func(t *testing.T) {
			testutil.RequireNoError(t, json.Unmarshal([]byte(`24`), &foo))

			testutil.RequireNoError(t, json.Unmarshal([]byte(`-24`), &foo))

			testutil.RequireEqual(t, true, foo.hasValue)
			testutil.RequireEqual(t, -24, foo.value)

			testutil.RequireError(t, foo.Validate())
			testutil.RequireEqual(t, false, foo.validated)

			testutil.RequirePanic(t, func() { foo.Must() })
			testutil.RequirePanic(t, func() { foo.Get() })
			testutil.RequirePanic(t, func() { foo.MarshalJSON() })
			testutil.RequirePanic(t, func() { foo.MarshalText() })
			testutil.RequirePanic(t, func() { foo.Value() })
		})
	})
}
