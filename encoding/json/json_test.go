package schemajson

import (
	"strings"
	"testing"

	"github.com/metafates/schema/internal/testutil"
	"github.com/metafates/schema/optional"
)

func TestJSON(t *testing.T) {
	type Mock struct {
		Foo string                 `json:"foo"`
		Bar optional.Positive[int] `json:"bar"`
	}

	for _, tc := range []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "valid json",
			json: `{"foo": "lorem ipsum", "bar": 249}`,
		},
		{
			name:    "invalid json",
			json:    `{"foo": lorem ipsum, ||||||| bar: 249}`,
			wantErr: true,
		},
		{
			name:    "validation error",
			json:    `{"foo": "lorem ipsum", "bar": -2}`,
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("unmarshal", func(t *testing.T) {
				var mock Mock
				err := Unmarshal([]byte(tc.json), &mock)

				if tc.wantErr {
					testutil.Error(t, err)
				} else {
					testutil.NoError(t, err)
				}
			})

			t.Run("decoder", func(t *testing.T) {
				var mock Mock

				err := NewDecoder(strings.NewReader(tc.json)).Decode(&mock)

				if tc.wantErr {
					testutil.Error(t, err)
				} else {
					testutil.NoError(t, err)
				}
			})
		})
	}
}
