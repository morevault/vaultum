// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"

	"github.com/morevault/vaultum/helper/identity"
)

func (c *Core) SendGroupUpdate(context.Context, *identity.Group) (bool, error) {
	return false, nil
}

func (c *Core) CreateEntity(ctx context.Context) (*identity.Entity, error) {
	return nil, nil
}
