// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package zos

import (
	hinfo "tailscale.com/hostinfo"
	"tailscale.com/tailcfg"
)

// HostInfo returns a partially populated Hostinfo for the current host.
func HostInfo() *tailcfg.Hostinfo {
	t := hinfo.New()
	t.IPNVersion = "zos-1.0.0"
	return t
}
