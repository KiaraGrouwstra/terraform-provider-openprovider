package domains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/charpand/terraform-provider-openprovider/internal/client"
)

// ListOptions narrows a domain listing. The zero value lists the account's
// domains as the API pages them by default.
type ListOptions struct {
	// FullName keeps only the domain with this full name (name and
	// extension), so a lookup does not page through the whole account.
	FullName string
	// Limit is the page size; the API's default applies when zero.
	Limit int
	// Offset is the number of domains to skip.
	Offset int
}

// ListWith retrieves a list of domains from the Openprovider API, narrowed by
// the given options.
//
// Endpoint: GET https://api.openprovider.eu/v1beta/domains
func ListWith(c *client.Client, opts ListOptions) ([]Domain, error) {
	query := url.Values{}
	if opts.FullName != "" {
		query.Set("full_name", opts.FullName)
	}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}

	path := "/v1beta/domains"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseURL, path), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var results ListDomainsResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results.Data.Results, nil
}
