package required

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestCustom_Parse(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   any
		wantErr bool
	}{
		{name: "valid int", value: 5},
		{name: "valid float", value: 5.2},
		{name: "invalid string", value: "hello", wantErr: true},
		{name: "invalid int", value: -2, wantErr: true},
		{name: "nil", value: nil, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var positive Positive[int]

			err := positive.Parse(tc.value)

			if tc.wantErr {
				testutil.Error(t, err)
				testutil.Equal(t, false, positive.validated)
				testutil.Equal(t, false, positive.hasValue)
				testutil.Equal(t, 0, positive.value)
				testutil.Panic(t, func() {
					positive.Get()
				})
			} else {
				testutil.NoError(t, err)
				testutil.Equal(t, true, positive.validated)
				testutil.Equal(t, true, positive.hasValue)
				testutil.Equal(t, reflect.ValueOf(tc.value).Convert(reflect.TypeFor[int]()).Interface().(int), positive.value)
				testutil.NoPanic(t, func() {
					positive.Get()
				})
			}
		})
	}
}

func TestRequired(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		testutil.NoError(t, json.Unmarshal([]byte(`null`), &foo))

		testutil.Equal(t, false, foo.hasValue)
		testutil.Equal(t, "", foo.value)

		testutil.Error(t, foo.TypeValidate())
		testutil.Equal(t, false, foo.validated)

		testutil.Panic(t, func() { foo.Get() })
		testutil.Panic(t, func() { foo.MarshalJSON() })
		testutil.Panic(t, func() { foo.MarshalText() })
		testutil.Panic(t, func() { foo.Value() })
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.NoError(t, json.Unmarshal([]byte(`-24`), &foo))

		testutil.Equal(t, true, foo.hasValue)
		testutil.Equal(t, -24, foo.value)

		testutil.Error(t, foo.TypeValidate())
		testutil.Equal(t, false, foo.validated)

		testutil.Panic(t, func() { foo.Get() })
		testutil.Panic(t, func() { foo.MarshalJSON() })
		testutil.Panic(t, func() { foo.MarshalText() })
		testutil.Panic(t, func() { foo.Value() })
	})

	t.Run("nested invalid value", func(t *testing.T) {
		type Foo struct {
			Field Positive[int]
		}

		var foo Any[Foo]

		testutil.NoError(t, json.Unmarshal([]byte(`{"field":-1}`), &foo))
		testutil.Error(t, foo.TypeValidate())
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.NoError(t, json.Unmarshal([]byte(`24`), &foo))

		testutil.Equal(t, true, foo.hasValue)
		testutil.Equal(t, 24, foo.value)

		testutil.NoError(t, foo.TypeValidate())
		testutil.Equal(t, true, foo.validated)

		testutil.NoPanic(t, func() { foo.Get() })
		testutil.NoPanic(t, func() { foo.MarshalJSON() })
		testutil.NoPanic(t, func() { foo.MarshalText() })
		testutil.NoPanic(t, func() { foo.Value() })

		t.Run("reuse as invalid", func(t *testing.T) {
			testutil.NoError(t, json.Unmarshal([]byte(`-24`), &foo))

			testutil.Equal(t, true, foo.hasValue)
			testutil.Equal(t, -24, foo.value)

			testutil.Error(t, foo.TypeValidate())
			testutil.Equal(t, false, foo.validated)

			testutil.Panic(t, func() { foo.Get() })
			testutil.Panic(t, func() { foo.MarshalJSON() })
			testutil.Panic(t, func() { foo.MarshalText() })
			testutil.Panic(t, func() { foo.Value() })
		})
	})
}
