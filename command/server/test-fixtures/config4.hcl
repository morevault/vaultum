
# SPDX-License-Identifier: MPL-2.0

disable_cache = true
ui = true

listener "tcp" {
	address = "127.0.0.1:8200"
}

storage "raft" {
	path = "/storage/path/raft"
	node_id = "raft1"
	performance_multiplier = 1
	foo = "bar"
	baz = true
}

cluster_addr = "127.0.0.1:8201"
