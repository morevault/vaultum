
# SPDX-License-Identifier: MPL-2.0

path "secret/foo" {
  policy = "write"
}

path "secret/bar/*" {
  capabilities = ["create", "read", "update"]
}
