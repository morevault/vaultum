// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"context"
	"os"

	"github.com/morevault/vaultum/api/v2"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/sys/unix"
)

func init() {
	getMmapFlags = getMmapFlagsLinux
}

func getMmapFlagsLinux(dbPath string) int {
	if api.ReadBaoVariable("BAO_RAFT_DISABLE_MAP_POPULATE") != "" {
		return 0
	}
	stat, err := os.Stat(dbPath)
	if err != nil {
		return 0
	}
	size := stat.Size()

	v, err := mem.VirtualMemoryWithContext(context.Background())
	if err != nil {
		return 0
	}

	// We won't worry about swap, since we already tell people not to use it.
	if v.Total > uint64(size) {
		return unix.MAP_POPULATE
	}

	return 0
}
