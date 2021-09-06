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

package zos

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
