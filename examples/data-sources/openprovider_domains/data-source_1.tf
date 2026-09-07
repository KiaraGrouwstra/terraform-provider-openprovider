# Every domain the account holds.
data "openprovider_domains" "all" {}

# One domain, or nothing: an empty list where the account does not hold it.
data "openprovider_domains" "example" {
  full_name = "example.com"
}

locals {
  held = length(data.openprovider_domains.example.domains) > 0
}
