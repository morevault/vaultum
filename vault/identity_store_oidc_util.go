// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"github.com/morevault/vaultum/helper/namespace"
)

func (i *IdentityStore) listNamespaces() []*namespace.Namespace {
	return []*namespace.Namespace{namespace.RootNamespace}
}
