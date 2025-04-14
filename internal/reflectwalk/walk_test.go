package reflectwalk

import (
	"reflect"
	"testing"
)

func TestWalkFields(t *testing.T) {
	type Foo struct{ Bar string }

	type Mock struct {
		Foo

		Name string
		Map  map[string][]int

		Anon struct {
			unexported bool
			Exported   uint8
		}

		Ptr  *Foo
		Ptr2 *Foo
	}

	mock := Mock{
		Foo:  Foo{Bar: "x"},
		Name: "x",
		Map: map[string][]int{
			"key":  {1, 2, 3},
			"key2": nil,
		},
		Anon: struct {
			unexported bool
			Exported   uint8
		}{
			unexported: true,
			Exported:   1,
		},
		Ptr: &Foo{
			Bar: "x",
		},
		Ptr2: nil,
	}

	want := map[string]any{
		"":               mock,
		".Anon":          mock.Anon,
		".Anon.Exported": mock.Anon.Exported,
		".Foo":           mock.Foo,
		".Foo.Bar":       mock.Foo.Bar,
		".Map":           mock.Map,
		".Map[key2]":     mock.Map["key2"],
		".Map[key]":      mock.Map["key"],
		".Map[key][0]":   mock.Map["key"][0],
		".Map[key][1]":   mock.Map["key"][1],
		".Map[key][2]":   mock.Map["key"][2],
		".Name":          mock.Name,
		".Ptr":           *mock.Ptr,
		".Ptr.Bar":       mock.Ptr.Bar,
		".Ptr2":          mock.Ptr2,
	}

	visited := make(map[string]any)

	err := WalkFields(mock, func(path string, value reflect.Value) error {
		visited[path] = value.Interface()

		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(want, visited) {
		t.Errorf("not equal:\nwant %#+v\ngot  %#+v", want, visited)
	}
}
