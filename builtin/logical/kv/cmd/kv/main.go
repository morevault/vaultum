package main

import (
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/morevault/vaultum/api/v2"
	"github.com/morevault/vaultum/builtin/logical/kv"
	"github.com/morevault/vaultum/sdk/v2/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: kv.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
