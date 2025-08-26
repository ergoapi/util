// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exnet

import "testing"

func TestGetIpsByDomain(t *testing.T) {
	ipaddr, err := GetIpsByDomain("blog.ysicing.net")
	if err != nil {
		t.Errorf("GetIpsByDomain error %v", err)
	}
	t.Logf("ips: %v", ipaddr)
}
