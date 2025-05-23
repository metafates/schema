package parse_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/metafates/schema/internal/testutil"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

func ExampleParse() {
	type User struct {
		Name    required.Any[string] `json:"such_tags_are_ignored_by_default"`
		Comment string
		Age     int
	}

	var user1, user2 User

	parse.Parse(map[string]any{
		"Name":         "john",
		"Comment":      "lorem ipsum",
		"Age":          99,
		"UnknownField": "this field will be ignored",
	}, &user1)

	parse.Parse(struct {
		Name, Comment string
		Age           uint8 // types will be converted
	}{
		Name:    "jane",
		Comment: "dolor sit",
		Age:     55,
	}, &user2)

	fmt.Printf("user1: name=%q comment=%q age=%v\n", user1.Name.Get(), user1.Comment, user1.Age)
	fmt.Printf("user2: name=%q comment=%q age=%v\n", user2.Name.Get(), user2.Comment, user2.Age)

	// Output:
	// user1: name="john" comment="lorem ipsum" age=99
	// user2: name="jane" comment="dolor sit" age=55
}

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
			{
				name:  "non-nil pointer",
				value: func() *int { n := 42; return &n }(),
				want:  42,
			},
			{
				name:  "nil pointer",
				value: (*int)(nil),
				want:  0,
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
			options []parse.Option
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
				want: Target{Name: sample.Name, Extra: Additional{CreatedAt: sample.Extra.CreatedAt}},
				value: map[string]any{
					"Name": sample.Name,
					"Extra": map[string]any{
						"CreatedAt": sample.Extra.CreatedAt,
					},
				},
			},
			{
				name: "invalid map with partial fields",
				want: Target{Bio: sample.Bio},
				value: map[string]any{
					"Bio": sample.Bio,
				},
				wantErr: true,
			},
			{
				name:    "invalid struct with partial fields",
				want:    Target{Name: sample.Name, Bio: sample.Bio},
				value:   struct{ Name, Bio string }{Name: sample.Name.Get(), Bio: sample.Bio},
				wantErr: true,
			},
			{
				name: "valid struct with partial fields",
				want: Target{Name: sample.Name, Extra: Additional{CreatedAt: sample.Extra.CreatedAt}},
				value: struct {
					Name  string
					Extra struct{ CreatedAt time.Time }
				}{
					Name:  sample.Name.Get(),
					Extra: struct{ CreatedAt time.Time }{CreatedAt: sample.Extra.CreatedAt.Get()},
				},
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
			{
				name: "unknown field",
				value: map[string]any{
					"Name": "john",
					"Foo":  9,
				},
				options: []parse.Option{parse.WithDisallowUnknownFields()},
				wantErr: true,
			},
			{
				name: "renamed fields",
				want: Target{Name: sample.Name, Extra: Additional{CreatedAt: sample.Extra.CreatedAt}},
				value: map[string]any{
					"name": sample.Name.Get(),
					"extra": map[string]any{
						"createdAt": sample.Extra.CreatedAt.Get(),
					},
				},
				options: []parse.Option{
					parse.WithDisallowUnknownFields(),
					parse.WithRenameFunc(func(s string) string {
						return strings.ToUpper(string(s[0])) + s[1:]
					}),
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var dst Target

				err := parse.Parse(tc.value, &dst, tc.options...)

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

//go:generate schemagen -type UserWithGen
type UserWithGen struct {
	ID    required.UUID[string]
	Name  required.Charset[string, charset.Print]
	Birth optional.InPast[time.Time]

	FavoriteNumber int

	Friends []Friend
}

func BenchmarkParse(b *testing.B) {
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

	b.Run("without codegen", func(b *testing.B) {
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
	})

	b.Run("with codegen", func(b *testing.B) {
		b.Run("parse", func(b *testing.B) {
			data := map[string]any{
				"ID":             "2c376d16-321d-43b3-8648-2e64798cc6b3",
				"Name":           "john",
				"FavoriteNumber": 42,
				"Friends": []map[string]any{
					{"ID": "7f735045-c8d2-4a60-9184-0fc033c40a6a", "Name": "jane"},
				},
			}

			var user UserWithGen

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

			var user UserWithGen

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
	})
}
