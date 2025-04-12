package testutil

import "testing"

func RequireEqual[T comparable](t *testing.T, want, actual T) {
	t.Helper()

	if want != actual {
		t.Fatalf("not equal: %v and %v", want, actual)
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
