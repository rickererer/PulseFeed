package interfaceshttpinteraction

import (
	"errors"
	"testing"

	domaininteraction "github.com/rickererer/PulseFeed/internal/domain/interaction"
)

func TestParseLimit(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    int
		wantErr error
	}{
		{"empty default", "", 0, nil},
		{"valid", "20", 20, nil},
		{"upper bound", "100", 100, nil},
		{"zero", "0", 0, domaininteraction.ErrInvalidLimit},
		{"negative", "-1", 0, domaininteraction.ErrInvalidLimit},
		{"over max", "101", 0, domaininteraction.ErrInvalidLimit},
		{"non numeric", "abc", 0, domaininteraction.ErrInvalidLimit},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLimit(tc.raw)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("parseLimit(%q) err = %v, want %v", tc.raw, err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("parseLimit(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}
