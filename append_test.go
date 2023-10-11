package multierror

import (
	"errors"
	"fmt"
	"testing"
)

func TestAppend_Error(t *testing.T) {
	original := &Error{
		Errors: []error{errors.New("foo")},
	}

	result := Append(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &Error{}
	result = Append(original, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test when a typed nil is passed
	var e *Error
	result = Append(e, errors.New("baz"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening
	original = &Error{
		Errors: []error{errors.New("foo")},
	}

	result = Append(original, Append(nil, errors.New("foo"), errors.New("bar")))
	if len(result.Errors) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilError(t *testing.T) {
	var err error
	result := Append(err, errors.New("bar"))
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *Error
	result := Append(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := Append(err, nilErr)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, errors.New("bar"))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, Append(nil, errors.New("bar")))
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_Errorf(t *testing.T) {
	var (
		A = errors.New("A")
		B = errors.New("B")
		E error
	)
	E = Append(
		Prefix(A, "some prefix"),
		fmt.Errorf("this is %w", B),
	)
	if !errors.Is(E, A) {
		t.Fatal("E is not A: ", E)
	}
	if !errors.Is(E, B) {
		t.Fatal("E is not B: ", B)
	}
	E = fmt.Errorf("error occurred: %w", E)
	if !errors.Is(E, A) {
		t.Fatal("E is not A: ", E)
	}
	if !errors.Is(E, B) {
		t.Fatal("E is not B: ", B)
	}
}
