
# SPDX-License-Identifier: MPL-2.0

disable_cache = true

backend "consul" {
    foo = "bar"
    disable_clustering = "true"
}

disable_clustering = false
