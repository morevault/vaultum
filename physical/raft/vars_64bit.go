
// SPDX-License-Identifier: MPL-2.0

//go:build !386 && !arm && !windows

package raft

const initialMmapSize = 100 * 1024 * 1024 * 1024 // 100GB
