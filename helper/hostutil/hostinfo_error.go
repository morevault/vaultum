// SPDX-License-Identifier: MPL-2.0

package hostutil

import "fmt"

// HostInfoError is a typed error for more convenient error checking.
type HostInfoError struct {
	Type string
	Err  error
}

func (e *HostInfoError) WrappedErrors() []error {
	return []error{e.Err}
}

func (e *HostInfoError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err.Error())
}
