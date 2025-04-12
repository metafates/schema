package testutil

import (
	"testing"
)

func RequireEqual[T comparable](t *testing.T, want, actual T) {
	t.Helper()

	if want != actual {
		t.Fatalf("not equal: want %v, got %v", want, actual)
	}
}

func RequireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func RequireError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("error is nil")
	}
}

func RequirePanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("did not panic")
		}
	}()

	f()
}

func RequireNoPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic")
		}
	}()

	f()
}
