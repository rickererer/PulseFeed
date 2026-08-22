package interfaceshttpfeed

import (
	"errors"
	"testing"

	domainfeed "github.com/rickererer/PulseFeed/internal/domain/feed"
)

func TestParseLimit(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    int
		wantErr error
	}{
		{"empty default", "", 0, nil},
		{"valid lower bound", "1", 1, nil},
		{"valid upper bound", "100", 100, nil},
		{"valid middle", "20", 20, nil},
		{"zero", "0", 0, domainfeed.ErrInvalidLimit},
		{"negative", "-5", 0, domainfeed.ErrInvalidLimit},
		{"over max", "101", 0, domainfeed.ErrInvalidLimit},
		{"way over max", "10000", 0, domainfeed.ErrInvalidLimit},
		{"non numeric", "abc", 0, domainfeed.ErrInvalidLimit},
		{"whitespace", "  10  ", 10, nil},
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
