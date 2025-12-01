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

func GetHostname() string {
	h := HostInfo()
	return h.Hostname
}

func GetOS() string {
	h := HostInfo()
	return h.OS
}

func GetDistro() string {
	h := HostInfo()
	return h.Distro
}

func GetDistroVersion() string {
	h := HostInfo()
	return h.DistroVersion
}

func GetDistroCodeName() string {
	h := HostInfo()
	return h.DistroCodeName
}

func GetOSVersion() string {
	h := HostInfo()
	return h.OSVersion
}

// IsDebian debian
func IsDebian() bool {
	h := HostInfo()
	return h.Distro == "debian"
}
