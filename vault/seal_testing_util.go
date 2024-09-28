// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"github.com/hashicorp/go-hclog"
	testing "github.com/mitchellh/go-testing-interface"
	"github.com/morevault/vaultum/sdk/v2/helper/logging"
	"github.com/morevault/vaultum/vault/seal"
	aeadwrapper "github.com/openbao/go-kms-wrapping/wrappers/aead/v2"
)

func NewTestSeal(t testing.T, opts *seal.TestSealOpts) Seal {
	t.Helper()
	if opts == nil {
		opts = &seal.TestSealOpts{}
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(hclog.Debug)
	}

	switch opts.StoredKeys {
	case seal.StoredKeysSupportedShamirRoot:
		newSeal := NewDefaultSeal(seal.NewAccess(aeadwrapper.NewShamirWrapper()))
		// Need StoredShares set or this will look like a legacy shamir seal.
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    1,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	case seal.StoredKeysNotSupported:
		newSeal := NewDefaultSeal(seal.NewAccess(aeadwrapper.NewShamirWrapper()))
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    0,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	default:
		access, _ := seal.NewTestSeal(opts)
		seal, err := NewAutoSeal(access)
		if err != nil {
			t.Fatal(err)
		}
		return seal
	}
}
