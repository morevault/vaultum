
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

template {
  destination = "/path/on/disk/where/template/will/render.txt"
  create_dest_dirs = true
  contents = "{{ keyOrDefault \"service/redis/maxconns@east-aws\" \"5\" }}"

  command = "restart service foo"
  command_timeout = "60s"

  error_on_missing_key = true
  perms = 0655
  backup = true
  left_delimiter  = "<<"
  right_delimiter = ">>"
  
  sandbox_path = "/path/on/disk/where"
  wait {
    min = "5s"
    max = "30s"
  }
  wait {
    min = "10s"
    max = "40s"
  }
}
