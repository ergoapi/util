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

package exstr

import "strings"

// Blacklist
func Blacklist(s string) bool {
	if strings.Contains(s, "<") {
		return true
	}

	if strings.Contains(s, ">") {
		return true
	}

	if strings.Contains(s, "&") {
		return true
	}

	if strings.Contains(s, "'") {
		return true
	}

	if strings.Contains(s, "\"") {
		return true
	}

	if strings.Contains(s, "file://") {
		return true
	}

	if strings.Contains(s, "../") {
		return true
	}

	if strings.Contains(s, "%") {
		return true
	}

	if strings.Contains(s, "=") {
		return true
	}

	if strings.Contains(s, "--") {
		return true
	}

	return false
}

// KubeBlacklist
func KubeBlacklist(s string) bool {
	if strings.HasPrefix(s, "kube-") {
		return true
	}

	if strings.HasSuffix(s, "-system") {
		return true
	}

	if strings.Contains(s, "cert-manager") {
		return true
	}

	if strings.Contains(s, "default") {
		return true
	}

	if strings.Contains(s, "observability") {
		return true
	}

	if strings.Contains(s, "tke-") {
		return true
	}

	if strings.Contains(s, "traefik") {
		return true
	}

	if strings.Contains(s, "velero") {
		return true
	}

	if strings.Contains(s, "ingress-nginx") {
		return true
	}

	return false
}
