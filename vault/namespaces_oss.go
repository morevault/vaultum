// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"

	"github.com/morevault/vaultum/helper/namespace"
)

func (c *Core) NamespaceByID(ctx context.Context, nsID string) (*namespace.Namespace, error) {
	return namespaceByID(ctx, nsID, c)
}

func (c *Core) ListNamespaces(includePath bool) []*namespace.Namespace {
	return []*namespace.Namespace{namespace.RootNamespace}
}
