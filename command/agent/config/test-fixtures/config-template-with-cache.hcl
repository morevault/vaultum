
# SPDX-License-Identifier: MPL-2.0

pid_file = "./pidfile"

auto_auth {
  method {
    type      = "aws"
    namespace = "/my-namespace"

    config = {
      role = "foobar"
    }
  }
}

cache {}

template {
  source      = "/path/on/disk/to/template.ctmpl"
  destination = "/path/on/disk/where/template/will/render.txt"
}
