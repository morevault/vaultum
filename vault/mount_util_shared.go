// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"github.com/morevault/vaultum/sdk/v2/logical"
)

func (c *Core) addBackendWriteForwardedPaths(backend logical.Backend, viewPath string) {
	paths := collectBackendSpecialPaths(backend, viewPath, func(specialPaths *logical.Paths) []string {
		return specialPaths.WriteForwardedStorage
	})

	c.logger.Trace("adding write forwarded paths", "paths", paths)
	c.writeForwardedPaths.AddPaths(paths)
}

func collectBackendSpecialPaths(backend logical.Backend, viewPath string, accessor func(specialPaths *logical.Paths) []string) []string {
	if backend == nil || backend.SpecialPaths() == nil {
		return nil
	}
	paths := accessor(backend.SpecialPaths())

	var ret []string
	for _, path := range paths {
		ret = append(ret, viewPath+path)
	}

	return ret
}
