
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/morevault/vaultum/plugins/database/influxdb"
	"github.com/morevault/vaultum/sdk/v2/database/dbplugin/v5"
)

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Run instantiates a Influxdb object, and runs the RPC server for the plugin
func Run() error {
	dbplugin.ServeMultiplex(influxdb.New)

	return nil
}
