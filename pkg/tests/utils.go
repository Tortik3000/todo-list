package tests

import (
	"errors"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nNot equal:\nexpected: %#v\nactual  : %#v", expected, actual)
	}
}

func RequireEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\nNot equal:\nexpected: %#v\nactual  : %#v", expected, actual)
	}
}

func RequireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Received unexpected error: %v", err)
	}
}

func RequireError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("An error is expected but got nil")
	}
}

func RequireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("Error expected to be '%v', but got '%v'", target, err)
	}
}
