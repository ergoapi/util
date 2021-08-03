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

package exhttp

// func GetTlsMsg(pemData string) string {
// 	block, rest := pem.Decode([]byte(pemData))
// 	if block == nil || len(rest) > 0 {
// 		zlog.Error("Certificate decoding error")
// 		return ""
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		zlog.Error("Certificate Parsing error, ", err.Error())
// 		return ""
// 	}
// 	result, err := certinfo.CertificateText(cert)
// 	if err != nil {
// 		zlog.Error("Certificate certinfo get error, ", err.Error())
// 		return ""
// 	}

// 	return result
// }

// func GetTlsMsgv2(pemData string) (*x509.Certificate, error) {
// 	block, _ := pem.Decode([]byte(pemData))
// 	if block == nil {
// 		zlog.Error("Certificate decoding error")
// 		return nil, errors.New("Certificate decoding error")
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		zlog.Error("Certificate Parsing error: ", err.Error())
// 		return nil, errors.New("Certificate Parsing error")
// 	}

// 	return cert, nil
// }
