package domains_test

import (
	"testing"

	"github.com/charpand/terraform-provider-openprovider/internal/client/domains"
)

func TestSplitFullName(t *testing.T) {
	cases := []struct {
		fullName      string
		wantName      string
		wantExtension string
		wantOK        bool
	}{
		{"example.com", "example", "com", true},
		// A second-level country suffix is one extension, not the last
		// label alone: splitting on the last dot would give "example.co"
		// and "uk" instead.
		{"example.co.uk", "example", "co.uk", true},
		{"no-dot", "", "", false},
		{"missing-name.", "", "", false},
		{".missing-extension", "", "", false},
	}
	for _, c := range cases {
		name, extension, ok := domains.SplitFullName(c.fullName)
		if ok != c.wantOK || name != c.wantName || extension != c.wantExtension {
			t.Errorf("SplitFullName(%q) = (%q, %q, %v), want (%q, %q, %v)",
				c.fullName, name, extension, ok, c.wantName, c.wantExtension, c.wantOK)
		}
	}
}
