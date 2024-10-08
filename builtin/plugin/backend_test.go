// SPDX-License-Identifier: MPL-2.0

package plugin_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/morevault/vaultum/api/v2"
	"github.com/morevault/vaultum/builtin/plugin"
	vaulthttp "github.com/morevault/vaultum/http"
	"github.com/morevault/vaultum/sdk/v2/helper/consts"
	"github.com/morevault/vaultum/sdk/v2/helper/logging"
	"github.com/morevault/vaultum/sdk/v2/helper/pluginutil"
	"github.com/morevault/vaultum/sdk/v2/logical"
	logicalPlugin "github.com/morevault/vaultum/sdk/v2/plugin"
	"github.com/morevault/vaultum/sdk/v2/plugin/mock"
	"github.com/morevault/vaultum/vault"
)

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = &plugin.PluginBackend{}
}

func TestBackend(t *testing.T) {
	pluginCmds := []string{"TestBackend_PluginMain", "TestBackend_PluginMain_Multiplexed"}

	for _, pluginCmd := range pluginCmds {
		t.Run(pluginCmd, func(t *testing.T) {
			config, cleanup := testConfig(t, pluginCmd)
			defer cleanup()

			_, err := plugin.Backend(context.Background(), config)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBackend_Factory(t *testing.T) {
	pluginCmds := []string{"TestBackend_PluginMain", "TestBackend_PluginMain_Multiplexed"}

	for _, pluginCmd := range pluginCmds {
		t.Run(pluginCmd, func(t *testing.T) {
			config, cleanup := testConfig(t, pluginCmd)
			defer cleanup()

			_, err := plugin.Factory(context.Background(), config)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBackend_PluginMain(t *testing.T) {
	args := []string{}
	if api.ReadBaoVariable(pluginutil.PluginUnwrapTokenEnv) == "" && api.ReadBaoVariable(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := api.ReadBaoVariable(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := logicalPlugin.Serve(&logicalPlugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_PluginMain_Multiplexed(t *testing.T) {
	args := []string{}
	if api.ReadBaoVariable(pluginutil.PluginUnwrapTokenEnv) == "" && api.ReadBaoVariable(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	caPEM := api.ReadBaoVariable(pluginutil.PluginCACertPEMEnv)
	if caPEM == "" {
		t.Fatal("CA cert not passed in")
	}

	args = append(args, fmt.Sprintf("--ca-cert=%s", caPEM))

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)
	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := logicalPlugin.ServeMultiplex(&logicalPlugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testConfig(t *testing.T, pluginCmd string) (*logical.BackendConfig, func()) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores

	core := cores[0]

	sys := vault.TestDynamicSystemView(core.Core, nil)

	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Debug),
		System: sys,
		Config: map[string]string{
			"plugin_name":    "mock-plugin",
			"plugin_type":    "secret",
			"plugin_version": "v0.0.0+mock",
		},
	}

	os.Setenv(pluginutil.PluginCACertPEMEnv, cluster.CACertPEMFile)

	vault.TestAddTestPlugin(t, core.Core, "mock-plugin", consts.PluginTypeSecrets, "", pluginCmd, []string{}, "")

	return config, func() {
		cluster.Cleanup()
	}
}
