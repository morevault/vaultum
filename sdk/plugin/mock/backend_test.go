// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"testing"

	"github.com/morevault/vaultum/sdk/v2/logical"
)

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = new(backend)
}
