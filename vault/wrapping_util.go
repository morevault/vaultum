
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"

	"github.com/morevault/vaultum/sdk/v2/logical"
)

func forwardWrapRequest(context.Context, *Core, *logical.Request, *logical.Response, *logical.Auth) (*logical.Response, error) {
	return nil, nil
}
