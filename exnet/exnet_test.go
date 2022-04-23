package exnet

import "testing"

func TestGetIpsByDomain(t *testing.T) {
	ipaddr, err := GetIpsByDomain("blog.ysicing.net")
	if err != nil {
		t.Errorf("GetIpsByDomain error %v", err)
	}
	t.Logf("ips: %v", ipaddr)
}
