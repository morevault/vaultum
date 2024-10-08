// SPDX-License-Identifier: MPL-2.0

package http

import (
	"net/http"

	"github.com/morevault/vaultum/vault"
)

func handleUnAuthenticatedInFlightRequest(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		currentInFlightReqMap := core.LoadInFlightReqData()

		respondOk(w, currentInFlightReqMap)
	})
}
