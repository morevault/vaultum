// SPDX-License-Identifier: MPL-2.0

package router

import (
	"testing"

	"github.com/morevault/vaultum/api/v2"
	"github.com/morevault/vaultum/builtin/credential/userpass"
	"github.com/morevault/vaultum/builtin/logical/pki"
	vaulthttp "github.com/morevault/vaultum/http"
	"github.com/morevault/vaultum/sdk/v2/logical"
	"github.com/morevault/vaultum/vault"
)

func TestRouter_MountSubpath_Checks(t *testing.T) {
	testRouter_MountSubpath(t, []string{"a/abcd/123", "abcd/123"})
	testRouter_MountSubpath(t, []string{"abcd/123", "a/abcd/123"})
	testRouter_MountSubpath(t, []string{"a/abcd/123", "abcd/123"})
}

func testRouter_MountSubpath(t *testing.T, mountPoints []string) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// Test auth
	authInput := &api.EnableAuthOptions{
		Type: "userpass",
	}
	for _, mp := range mountPoints {
		t.Logf("mounting %s", "auth/"+mp)
		var err error
		err = client.Sys().EnableAuthWithOptions("auth/"+mp, authInput)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Test secrets
	mountInput := &api.MountInput{
		Type: "pki",
	}
	for _, mp := range mountPoints {
		t.Logf("mounting %s", "s/"+mp)
		var err error
		err = client.Sys().Mount("s/"+mp, mountInput)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)
	t.Logf("Done: %#v", mountPoints)
}

func TestRouter_UnmountRollbackIsntFatal(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"noop": vault.NoopBackendRollbackErrFactory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	if err := client.Sys().Mount("noop", &api.MountInput{
		Type: "noop",
	}); err != nil {
		t.Fatalf("failed to mount PKI: %v", err)
	}

	if _, err := client.Logical().Write("sys/plugins/reload/backend", map[string]interface{}{
		"mounts": "noop",
	}); err != nil {
		t.Fatalf("expected reload of noop with broken periodic func to succeed; got err=%v", err)
	}

	if _, err := client.Logical().Write("sys/remount", map[string]interface{}{
		"from": "noop",
		"to":   "noop-to",
	}); err != nil {
		t.Fatalf("expected remount of noop with broken periodic func to succeed; got err=%v", err)
	}

	cluster.EnsureCoresSealed(t)
	cluster.UnsealCores(t)
}
