// SPDX-License-Identifier: MPL-2.0

package totp

import (
	"context"
	"strings"
	"time"

	"github.com/morevault/vaultum/sdk/v2/framework"
	"github.com/morevault/vaultum/sdk/v2/logical"
	cache "github.com/patrickmn/go-cache"
)

const operationPrefixTOTP = "totp"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"key/",
			},
		},

		Paths: []*framework.Path{
			pathListKeys(&b),
			pathKeys(&b),
			pathCode(&b),
		},

		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	b.usedCodes = cache.New(0, 30*time.Second)

	return &b
}

type backend struct {
	*framework.Backend

	usedCodes *cache.Cache
}

const backendHelp = `
The TOTP backend dynamically generates time-based one-time use passwords.
`
