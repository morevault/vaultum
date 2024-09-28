// SPDX-License-Identifier: MPL-2.0

package main // import "github.com/morevault/vault"

import (
	"os"

	"github.com/morevault/vaultum/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
