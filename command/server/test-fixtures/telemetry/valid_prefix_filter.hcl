
# SPDX-License-Identifier: MPL-2.0

ui            = true

telemetry {
  statsd_address = "foo"
  prefix_filter = ["-vault.expire", "-vault.audit", "+vault.expire.num_irrevocable_leases"]
}