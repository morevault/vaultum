// SPDX-License-Identifier: MPL-2.0

package http

import (
	"fmt"
	"net/http"

	"github.com/morevault/vaultum/helper/metricsutil"
	"github.com/morevault/vaultum/sdk/v2/logical"
	"github.com/morevault/vaultum/vault"
)

func handleMetricsUnauthenticated(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &logical.Request{Headers: r.Header}

		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		format := r.Form.Get("format")
		if format == "" {
			format = metricsutil.FormatFromRequest(req)
		}

		// Define response
		resp := core.MetricsHelper().ResponseForFormat(format)

		// Manually extract the logical response and send back the information
		status := resp.Data[logical.HTTPStatusCode].(int)
		w.Header().Set("Content-Type", resp.Data[logical.HTTPContentType].(string))
		switch v := resp.Data[logical.HTTPRawBody].(type) {
		case string:
			w.WriteHeader(status)
			w.Write([]byte(v))
		case []byte:
			w.WriteHeader(status)
			w.Write(v)
		default:
			respondError(w, http.StatusInternalServerError, fmt.Errorf("wrong response returned"))
		}
	})
}
