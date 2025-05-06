package typeconv

import (
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "empty string",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "single valid tag",
			input: `json:"name"`,
			want:  map[string]string{"json": "name"},
		},
		{
			name:  "multiple tags with spaces",
			input: `  json:"name"   xml:"id"  `,
			want:  map[string]string{"json": "name", "xml": "id"},
		},
		{
			name:  "escaped quote in value",
			input: `json:"na\"me"`,
			want:  map[string]string{"json": `na\"me`},
		},
		{
			name:  "unescaped quote in value",
			input: `json:"na"me"`,
			want:  map[string]string{"json": "na"},
		},
		{
			name:  "missing closing quote",
			input: `json:"name`,
			want:  map[string]string{},
		},
		{
			name:  "key with no colon",
			input: `json`,
			want:  map[string]string{},
		},
		{
			name:  "key with colon but no quote",
			input: `json:name`,
			want:  map[string]string{},
		},
		{
			name:  "empty key",
			input: `:"value"`,
			want:  map[string]string{"": "value"},
		},
		{
			name:  "duplicate keys",
			input: `json:"a" json:"b"`,
			want:  map[string]string{"json": "b"},
		},
		{
			name:  "key with trailing spaces and valid value",
			input: `  json   :"name"  `,
			want:  map[string]string{"json   ": "name"},
		},
		{
			name:  "value with backslashes",
			input: `json:"a\\b"`,
			want:  map[string]string{"json": `a\\b`},
		},
		{
			name:  "empty value",
			input: `json:""`,
			want:  map[string]string{"json": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseStructTags(tt.input)

			testutil.DeepEqual(t, tt.want, actual)
		})
	}
}
