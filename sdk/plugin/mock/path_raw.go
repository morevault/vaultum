
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"

	"github.com/morevault/vaultum/sdk/v2/framework"
	"github.com/morevault/vaultum/sdk/v2/logical"
)

// pathRaw is used to test raw responses.
func pathRaw(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "raw",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRawRead,
		},
	}
}

func (b *backend) pathRawRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "text/plain",
			logical.HTTPRawBody:     []byte("Response"),
			logical.HTTPStatusCode:  200,
		},
	}, nil
}
