// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"

	"github.com/morevault/vaultum/helper/namespace"
)

var NamespaceByID func(context.Context, string, *Core) (*namespace.Namespace, error) = namespaceByID

func namespaceByID(ctx context.Context, nsID string, c *Core) (*namespace.Namespace, error) {
	if nsID == namespace.RootNamespaceID {
		return namespace.RootNamespace, nil
	}
	return nil, namespace.ErrNoNamespace
}
