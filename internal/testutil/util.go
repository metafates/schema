package testutil

import "reflect"

type Handle interface {
	Helper()
	Fatalf(format string, args ...any)
}

func Equal[T comparable](t Handle, want, actual T) {
	t.Helper()

	if want != actual {
		t.Fatalf("not equal:\nwant %+v\ngot  %+v", want, actual)
	}
}

func DeepEqual(t Handle, want, actual any) {
	t.Helper()

	if !reflect.DeepEqual(want, actual) {
		t.Fatalf("not equal:\nwant %+v\ngot  %+v", want, actual)
	}
}

func NoError(t Handle, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Error(t Handle, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("error is nil")
	}
}

func Panic(t Handle, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("did not panic")
		}
	}()

	f()
}

func NoPanic(t Handle, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic")
		}
	}()

	f()
}
