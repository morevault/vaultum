// SPDX-License-Identifier: MPL-2.0

//go:build !ui

package http

import (
	"net/http"
)

func init() {
	// set uiBuiltIn to false to indicate the ui is not built in. See
	// http/handler.go
	uiBuiltIn = false
}

// assetFS is a stub for building Vault without a UI.
func assetFS() http.FileSystem {
	return nil
}
