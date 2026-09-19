// Package domains_test contains tests for the domains package.
package domains_test

import (
	"testing"

	"github.com/charpand/terraform-provider-openprovider/internal/client/domains"
	"github.com/charpand/terraform-provider-openprovider/internal/testutils"
)

func TestCheckDomains(t *testing.T) {
	apiClient := testutils.SetupTestClient()

	results, err := domains.Check(apiClient, []domains.CheckDomain{{Name: "example", Extension: "com"}})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) == 0 {
		t.Log("Note: No check results returned by mock server (check your swagger examples)")
	}
}
