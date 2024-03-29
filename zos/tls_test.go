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

import "testing"

func TestExTLSCheck(t *testing.T) {
	const certPem = `-----BEGIN CERTIFICATE-----
	MIIEZjCCAs6gAwIBAgIQbzMRx5UjegIWsBoT9p2EKDANBgkqhkiG9w0BAQsFADCB
	gTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSswKQYDVQQLDCJ5c2lj
	aW5nQFlzaUNJbmdkZU1hY0Jvb2stUHJvLmxvY2FsMTIwMAYDVQQDDClta2NlcnQg
	eXNpY2luZ0BZc2lDSW5nZGVNYWNCb29rLVByby5sb2NhbDAeFw0xOTA2MDEwMDAw
	MDBaFw0zMDA4MjQwNzM3MjBaMFExJzAlBgNVBAoTHm1rY2VydCBkZXZlbG9wbWVu
	dCBjZXJ0aWZpY2F0ZTEmMCQGA1UECwwdeXNpY2luZ0Bib2dvbiAoWXNpQ0luZyBa
	aGVuZykwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC8yzPIrw5l8rKb
	Re+N4x9ItBlnsYffOwXijfbrm6kphCKzfjFqJMk0QGmFjfiaP7iRZ1UrucY+lgOH
	p5ioKU8cUvryXxIKUBVof7X26UZg8YrfDdzXmDJTWRGEqxfRKjQWkNvAMy43oHno
	D/gEv8dStIOtizRlRzLd6Ax7PCS2n8YeUKSALgsGASTw+wEIjfaZVRpvee0r1zxg
	vpfuRLUYwweS/ZxupwMiBznX3tq70wePXn6qd0Ax0KJoJX+HLR1pYRl1DMt9MOOV
	WJp00cywiPmttcUaViwClgQn8UZRAsPdsyDj5F1nKo0iLD6Ta9FyDZXRBxP/6PVO
	sN3Gvev1AgMBAAGjgYgwgYUwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsG
	AQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUFixKysq+hwVH/Wor5WmA
	P8+uenQwLwYDVR0RBCgwJoINd3d3LmJhaWR1LmNvbYIJbG9jYWxob3N0hwR/AAAB
	hwSsEEhYMA0GCSqGSIb3DQEBCwUAA4IBgQDDHeRtyfwXzRJnfiICxmWf/aQTRIZB
	WzQqrHSIpTiUrZpB5UEMLDmiy16NjQKy+4261Bzh0N/T/RXPtiW2ktoVwuoELVrb
	a+0idjlvT3xPr8+aWc/PNCq/D6vkMrDYU6L/apfHgOhbO6Wxw/JjE8c1qUrtePDT
	bRHx8kme7hoi7t+4Fk5hwxMI/4wvkAgYvn3qRCHZuSNTgNyv9V5/WGXflUXOp4Di
	Bp+fYtidH4Yix6ESYVxmQeElP59bIJtzDUjIC4DcTcug2al99nUp6B5ZgAWTw5jX
	TtgK+tqbMiaoyITslD+cJjB7V2fpc6Ltkf9wf3QPgGIoeQqfn8RtP+4hyqrYlC9k
	5msEe0Ros1qH/1nEWirP60Q3TRBUWxVOIR4BTGMUOlYf+C/J+xbsJIBb6Sj+LhiG
	38Yn8Rn7FgSYXuqeaGIXvsncrMJQJ77OJIDL2TLYqLIJNMw/3ogX5SJLX6tduBWB
	ciQJgeNl87QtzMHspX5bAFBIORSpCVVnrFQ=
	-----END CERTIFICATE-----`

	tlsres, err := TLSCheck(certPem)
	if err == nil {
		t.Logf("domain %v", tlsres.DNSNames)
	}
}
