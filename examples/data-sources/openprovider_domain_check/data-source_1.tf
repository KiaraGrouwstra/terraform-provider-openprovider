data "openprovider_domain_check" "example" {
  domain = "example.com"
}

# Register the domain only where nobody holds it yet.
resource "openprovider_domain" "example" {
  count = data.openprovider_domain_check.example.available ? 1 : 0

  domain       = "example.com"
  owner_handle = "XX000000-NL"
}
