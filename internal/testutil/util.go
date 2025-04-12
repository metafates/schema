package testutil

import (
	"testing"
)

func Equal[T comparable](t *testing.T, want, actual T) {
	t.Helper()

	if want != actual {
		t.Fatalf("not equal: want %v, got %v", want, actual)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("error is nil")
	}
}

func Panic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("did not panic")
		}
	}()

	f()
}

func NoPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic")
		}
	}()

	f()
}
