package parse_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/metafates/schema/internal/testutil"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

func TestParse(t *testing.T) {
	t.Run("basic dst", func(t *testing.T) {
		for _, tc := range []struct {
			name    string
			value   any
			want    int
			wantErr bool
		}{
			{
				name:  "same type",
				value: int(42),
				want:  42,
			},
			{
				name:  "type conversion",
				value: float64(42.000000),
				want:  42,
			},
			{
				name:    "invalid type",
				value:   "hello",
				want:    42,
				wantErr: true,
			},
			{
				name:  "loosely type conversion",
				value: 42.42,
				want:  42,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var dst int

				err := parse.Parse(tc.value, &dst)

				if tc.wantErr {
					testutil.Error(t, err)
				} else {
					testutil.NoError(t, err)
					testutil.Equal(t, tc.want, dst)
				}
			})
		}
	})

	t.Run("slice dst", func(t *testing.T) {
		type Foo struct {
			Name string
		}

		for _, tc := range []struct {
			name    string
			want    []Foo
			value   any
			wantErr bool
		}{
			{
				name: "valid maps",
				want: []Foo{
					{Name: "Foo"},
					{Name: "Bar"},
				},
				value: []map[string]string{
					{"Name": "Foo"},
					{"Name": "Bar"},
				},
			},
			{
				name: "valid structs",
				want: []Foo{
					{Name: "Foo"},
					{Name: "Bar"},
				},
				value: []struct{ Name string }{
					{Name: "Foo"},
					{Name: "Bar"},
				},
			},
			{
				name:  "nil",
				want:  nil,
				value: nil,
			},
			{
				name: "partial",
				want: []Foo{
					{Name: ""},
					{Name: "Bar"},
				},
				value: []map[string]string{
					{"Surname": "Foo"},
					{"Name": "Bar"},
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var dst []Foo

				err := parse.Parse(tc.value, &dst)

				if tc.wantErr {
					testutil.Error(t, err)
				} else {
					testutil.NoError(t, err)
					testutil.DeepEqual(t, tc.want, dst)
				}
			})
		}
	})

	t.Run("struct dst", func(t *testing.T) {
		type Additional struct {
			CreatedAt required.InPast[time.Time]
			IsAdmin   bool
			Theme     optional.Any[string]
			Foo       optional.Any[int]
		}

		type Target struct {
			Name required.NonZero[string] `json:"json_tag_should_be_ignored"`
			Bio  string
			IDs  []int

			Extra Additional
		}

		sample := Target{
			Name: func() (name required.NonZero[string]) {
				name.MustParse("john")

				return name
			}(),
			IDs: []int{1, 100},
			Bio: "...",
			Extra: Additional{
				CreatedAt: func() (v required.InPast[time.Time]) {
					v.MustParse(time.Now().Add(-time.Hour))

					return v
				}(),
				IsAdmin: true,
				Theme: func() (v optional.Any[string]) {
					v.MustParse("dark")

					return v
				}(),
			},
		}

		for _, tc := range []struct {
			name    string
			value   any
			want    Target
			wantErr bool
		}{
			{
				name: "valid map with all fields",
				want: sample,
				value: map[string]any{
					"Name": sample.Name.Get(),
					"IDs":  []uint8{1, 100}, // types are converted
					"Bio":  sample.Bio,
					"Extra": map[string]any{
						"CreatedAt": sample.Extra.CreatedAt.Get(),
						"IsAdmin":   sample.Extra.IsAdmin,
						"Theme":     sample.Extra.Theme.GetPtr(),
						"Foo":       sample.Extra.Foo.GetPtr(),
					},
				},
			},
			{
				name: "valid map with partial fields",
				want: Target{Bio: sample.Bio},
				value: map[string]any{
					"Bio": sample.Bio,
				},
			},
			{
				name:  "valid struct with partial fields",
				want:  Target{Name: sample.Name, Bio: sample.Bio},
				value: struct{ Name, Bio string }{Name: sample.Name.Get(), Bio: sample.Bio},
			},
			{
				name:  "valid same struct",
				want:  sample,
				value: sample,
			},
			{
				name: "invalid type",
				want: sample,
				value: map[string]any{
					"Name": 42.42,
				},
				wantErr: true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var dst Target

				err := parse.Parse(tc.value, &dst)

				if tc.wantErr {
					testutil.Error(t, err)
				} else {
					testutil.NoError(t, err)
					testutil.DeepEqual(t, tc.want, dst)
				}
			})
		}
	})
}

func BenchmarkParse(b *testing.B) {
	type Friend struct {
		ID   required.UUID[string]
		Name required.Charset[string, charset.Print]
	}

	type User struct {
		ID    required.UUID[string]
		Name  required.Charset[string, charset.Print]
		Birth optional.InPast[time.Time]

		FavoriteNumber int

		Friends []Friend
	}

	b.Run("manual", func(b *testing.B) {
		var user User

		for b.Loop() {
			if err := user.ID.Parse("2c376d16-321d-43b3-8648-2e64798cc6b3"); err != nil {
				b.Fatal(err)
			}

			if err := user.Name.Parse("john"); err != nil {
				b.Fatal(err)
			}

			user.FavoriteNumber = 42
			user.Friends = make([]Friend, 1)

			if err := user.Friends[0].ID.Parse("7f735045-c8d2-4a60-9184-0fc033c40a6a"); err != nil {
				b.Fatal(err)
			}

			if err := user.Friends[0].Name.Parse("jane"); err != nil {
				b.Fatal(err)
			}
		}

		_ = user
	})

	b.Run("parse", func(b *testing.B) {
		data := map[string]any{
			"ID":             "2c376d16-321d-43b3-8648-2e64798cc6b3",
			"Name":           "john",
			"FavoriteNumber": 42,
			"Friends": []map[string]any{
				{"ID": "7f735045-c8d2-4a60-9184-0fc033c40a6a", "Name": "jane"},
			},
		}

		var user User

		for b.Loop() {
			if err := parse.Parse(data, &user); err != nil {
				b.Fatal(err)
			}
		}

		_ = user
	})

	b.Run("unmarshal", func(b *testing.B) {
		data := []byte(`
		{
			"ID": "2c376d16-321d-43b3-8648-2e64798cc6b3",
			"Name": "john",
			"FavoriteNumber": 42,
			"Friends": [
				{"ID": "7f735045-c8d2-4a60-9184-0fc033c40a6a", "Name": "jane"}
			]
		}
		`)

		var user User

		for b.Loop() {
			if err := json.Unmarshal(data, &user); err != nil {
				b.Fatal(err)
			}

			if err := validate.Validate(&user); err != nil {
				b.Fatal(err)
			}
		}

		_ = user
	})
}
