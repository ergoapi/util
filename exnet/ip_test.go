//  Copyright (c) 2021. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package exnet

import (
	"net"
	"testing"
)

func TestIsPrivateIP(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1.1.1.1",
			args: args{
				ip: net.ParseIP("1.1.1.1"),
			},
			want: false,
		},
		{
			name: "10.101.101.10",
			args: args{
				ip: net.ParseIP("10.101.101.10"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivateIP(tt.args.ip); got != tt.want {
				t.Errorf("IsPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatMacAddr(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"92:10:31:a3:ec:37", "92:10:31:a3:ec:37"},
		{"92-10-31-A3-EC-37", "92:10:31:a3:ec:37"},
	}

	for _, c := range cases {
		if FormatMacAddr(c.in) != c.out {
			t.Errorf(" %s => %s != %s", c.in, FormatMacAddr(c.in), c.out)
		}
	}
}

func TestIP2Number(t *testing.T) {
	for _, addr := range []string{"192.168.23.1", "255.255.255.255", "0.0.0.0"} {
		num, err := IP2Number(addr)
		if err != nil {
			t.Errorf("IP2Number error %s %s", addr, err)
		}
		ipstr := Number2IP(num)
		if ipstr != addr {
			t.Errorf("%s != %s", ipstr, addr)
		}
	}
}

func TestIPV4Addr_StepDown(t *testing.T) {
	ipaddr, err := NewIPV4Addr("192.168.222.253")
	if err != nil {
		t.Errorf("NewIPV4Addr error %s", err)
	}
	t.Logf("Network: %s Broadcast: %s Client: %s", ipaddr.NetAddr(24), ipaddr.BroadcastAddr(24), ipaddr.CliAddr(24))
	t.Logf("Stepup: %s", ipaddr.StepUp())
	t.Logf("Stepdown: %s", ipaddr.StepDown())
}
