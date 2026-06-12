package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapPreservesCause(t *testing.T) {
	cause := errors.New("connection refused")
	err := Wrap(KindUnavailable, "storage.Connect", cause)

	if !errors.Is(err, cause) {
		t.Fatal("wrapped error must match its cause via errors.Is")
	}
	if KindOf(err) != KindUnavailable {
		t.Fatalf("KindOf = %s, want %s", KindOf(err), KindUnavailable)
	}
}

func TestWrapNilReturnsNil(t *testing.T) {
	if err := Wrap(KindInternal, "op", nil); err != nil {
		t.Fatalf("Wrap(nil) = %v, want nil", err)
	}
	if err := Wrapf(KindInternal, "op", nil, "msg"); err != nil {
		t.Fatalf("Wrapf(nil) = %v, want nil", err)
	}
}

func TestKindOfUnclassifiedErrorIsInternal(t *testing.T) {
	if kind := KindOf(errors.New("plain")); kind != KindInternal {
		t.Fatalf("KindOf = %s, want %s", kind, KindInternal)
	}
}

func TestKindOfFindsErrorDeepInChain(t *testing.T) {
	err := fmt.Errorf("handler: %w", Wrap(KindConflict, "incident.Close", ErrIncidentAlreadyClosed))

	if !Is(err, KindConflict) {
		t.Fatal("expected KindConflict in wrapped chain")
	}
	if !errors.Is(err, ErrIncidentAlreadyClosed) {
		t.Fatal("expected sentinel ErrIncidentAlreadyClosed in chain")
	}
}

func TestErrorMessageFormat(t *testing.T) {
	cases := []struct {
		name string
		err  *Error
		want string
	}{
		{"message only", New(KindValidation, "incident.Create", "title is empty"), "incident.Create: title is empty"},
		{"cause only", Wrap(KindInternal, "config.Load", errors.New("boom")), "config.Load: boom"},
		{"message and cause", Wrapf(KindInternal, "config.Load", errors.New("boom"), "read env"), "config.Load: read env: boom"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}
