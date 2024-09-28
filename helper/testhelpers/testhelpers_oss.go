// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"github.com/mitchellh/go-testing-interface"
	"github.com/morevault/vaultum/vault"
)

// WaitForActiveNodeAndStandbys does nothing more than wait for the active node
// on OSS. On enterprise it waits for perf standbys to be healthy too.
func WaitForActiveNodeAndStandbys(t testing.T, cluster *vault.TestCluster) {
	WaitForActiveNode(t, cluster)
	for _, core := range cluster.Cores {
		if standby, _ := core.Core.Standby(); standby {
			WaitForStandbyNode(t, core)
		}
	}
}
