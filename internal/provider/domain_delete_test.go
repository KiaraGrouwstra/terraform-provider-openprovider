package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charpand/terraform-provider-openprovider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// domainAPI stands in for OpenProvider with one domain in the account,
// example.com as id 123, and records the path of every delete it is sent.
func domainAPI(t *testing.T, deleted *[]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1beta/domains":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"data": map[string]any{
					"results": []map[string]any{{
						"id":     123,
						"domain": map[string]string{"name": "example", "extension": "com"},
						"status": "ACT",
					}},
					"total": 1,
				},
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/v1beta/domains/123":
			*deleted = append(*deleted, r.URL.Path)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"data": map[string]any{"success": true},
			})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// domainState is the state of example.com as this provider would hold it,
// with `on_destroy` as given (null for a state written before it existed).
func domainState(t *testing.T, r *DomainResource, onDestroy types.String) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema}
	diags := state.Set(ctx, DomainModel{
		ID:          types.StringValue("example.com"),
		Domain:      types.StringValue("example.com"),
		OwnerHandle: types.StringValue("XX000000-NL"),
		OnDestroy:   onDestroy,
		DnssecKeys:  types.ListNull(types.ObjectType{AttrTypes: dnssecKeysAttrTypes}),
	})
	if diags.HasError() {
		t.Fatalf("setting state: %v", diags)
	}
	return state
}

func deleteDomain(t *testing.T, onDestroy types.String) (deleted []string, errored bool) {
	t.Helper()
	server := domainAPI(t, &deleted)
	defer server.Close()
	r := &DomainResource{client: client.NewClient(client.Config{
		BaseURL:    server.URL,
		Token:      "test",
		HTTPClient: server.Client(),
	})}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), resource.DeleteRequest{State: domainState(t, r, onDestroy)}, resp)
	return deleted, resp.Diagnostics.HasError()
}

func TestDomainResourceDeleteRetains(t *testing.T) {
	deleted, errored := deleteDomain(t, types.StringValue(onDestroyRetain))
	if errored {
		t.Fatal("retaining a domain should not error")
	}
	if len(deleted) != 0 {
		t.Fatalf("retaining a domain should send no delete, sent %v", deleted)
	}
}

func TestDomainResourceDeleteRetainsWhenUnset(t *testing.T) {
	deleted, errored := deleteDomain(t, types.StringNull())
	if errored {
		t.Fatal("a state without on_destroy should destroy without error")
	}
	if len(deleted) != 0 {
		t.Fatalf("a state without on_destroy should send no delete, sent %v", deleted)
	}
}

func TestDomainResourceDeleteDeletes(t *testing.T) {
	deleted, errored := deleteDomain(t, types.StringValue(onDestroyDelete))
	if errored {
		t.Fatal("deleting a domain should not error")
	}
	if len(deleted) != 1 || deleted[0] != "/v1beta/domains/123" {
		t.Fatalf("deleting a domain should send one delete for its id, sent %v", deleted)
	}
}
