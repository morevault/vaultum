
# SPDX-License-Identifier: MPL-2.0

pid_file = "./pidfile"

auto_auth {
	method {
		type = "aws"
		wrap_ttl = 300
		config = {
			role = "foobar"
		}
	}

	sink {
		type = "file"
		wrap_ttl = 300
		config = {
			path = "/tmp/file-foo"
		}
	}
}
