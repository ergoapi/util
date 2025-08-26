// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package validation provides validation utilities.
package validation

import (
	"net"
	"regexp"
	"strings"
)

var FunctionReg *regexp.Regexp
var UUIDReg *regexp.Regexp
var UUIDExactReg *regexp.Regexp
var IntegerReg *regexp.Regexp
var FloatReg *regexp.Regexp
var MacaddrReg *regexp.Regexp
var CompactMacaddrReg *regexp.Regexp
var NsptrReg *regexp.Regexp
var NameReg *regexp.Regexp
var DomainnameReg *regexp.Regexp
var DomainsrvReg *regexp.Regexp
var SizeReg *regexp.Regexp
var MonthReg *regexp.Regexp
var DateReg *regexp.Regexp
var DateCompactReg *regexp.Regexp
var IsoTimeReg *regexp.Regexp
var IsoNoSecondTimeReg *regexp.Regexp
var FullisoTimeReg *regexp.Regexp
var IsoTimeReg2 *regexp.Regexp
var IsoNoSecondTimeReg2 *regexp.Regexp
var FullisoTimeReg2 *regexp.Regexp
var ZstackTimeReg *regexp.Regexp
var CompactTimeReg *regexp.Regexp
var MysqlTimeReg *regexp.Regexp
var NormalTimeReg *regexp.Regexp
var FullnormalTimeReg *regexp.Regexp
var Rfc2882TimeReg *regexp.Regexp
var EmailReg *regexp.Regexp
var ChinaMobileReg *regexp.Regexp
var FsFormatReg *regexp.Regexp
var UsCurrencyReg *regexp.Regexp
var EuCurrencyReg *regexp.Regexp

func init() {
	FunctionReg = regexp.MustCompile(`^\w+\(.*\)$`)
	UUIDReg = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)
	UUIDExactReg = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	IntegerReg = regexp.MustCompile(`^[0-9]+$`)
	FloatReg = regexp.MustCompile(`^\d+(\.\d*)?$`)
	MacaddrReg = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
	CompactMacaddrReg = regexp.MustCompile(`^[0-9a-fA-F]{12}$`)
	NsptrReg = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\.in-addr\.arpa$`)
	NameReg = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._@-]*$`)
	DomainnameReg = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)
	SizeReg = regexp.MustCompile(`^\d+[bBkKmMgG]?$`)
	MonthReg = regexp.MustCompile(`^\d{4}-\d{2}$`)
	DateReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	DateCompactReg = regexp.MustCompile(`^\d{8}$`)
	IsoTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)
	IsoNoSecondTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)
	FullisoTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3,9}(Z|[+-]\d{2}:\d{2})$`)
	IsoTimeReg2 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)
	IsoNoSecondTimeReg2 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)
	FullisoTimeReg2 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3,9}(Z|[+-]\d{2}:\d{2})$`)
	CompactTimeReg = regexp.MustCompile(`^\d{14}$`)
	ZstackTimeReg = regexp.MustCompile(`^\w+ \d{1,2}, \d{4} \d{1,2}:\d{1,2}:\d{1,2} (AM|PM)$`) //ZStack time format "Apr 1, 2019 3:23:17 PM"
	MysqlTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	NormalTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`)
	FullnormalTimeReg = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{6}$`)
	Rfc2882TimeReg = regexp.MustCompile(`[A-Z][a-z]{2}, [0-9]{1,2} [A-Z][a-z]{2} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Z]{3}`)
	EmailReg = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,4}$`)
	ChinaMobileReg = regexp.MustCompile(`^1[0-9-]{10}$`)
	FsFormatReg = regexp.MustCompile(`^(ext|fat|hfs|xfs|swap|ntfs|reiserfs|ufs|btrfs)`)
	UsCurrencyReg = regexp.MustCompile(`^[+-]?(\d{0,3}|((\d{1,3},)+\d{3}))(\.\d*)?$`)
	EuCurrencyReg = regexp.MustCompile(`^[+-]?(\d{0,3}|((\d{1,3}\.)+\d{3}))(,\d*)?$`)
}

func MatchFunction(str string) bool {
	return FunctionReg.MatchString(str)
}

func MatchUUID(str string) bool {
	return UUIDReg.MatchString(str)
}

func MatchUUIDExact(str string) bool {
	return UUIDExactReg.MatchString(str)
}

func MatchInteger(str string) bool {
	return IntegerReg.MatchString(str)
}

func MatchFloat(str string) bool {
	return FloatReg.MatchString(str)
}

func MatchMacAddr(str string) bool {
	return MacaddrReg.MatchString(str)
}

func MatchCompactMacAddr(str string) bool {
	return CompactMacaddrReg.MatchString(str)
}

func MatchIP4Addr(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && !strings.Contains(str, ":")
}

func MatchCIDR(str string) bool {
	ip, _, err := net.ParseCIDR(str)
	if err != nil {
		return false
	}
	return ip != nil && !strings.Contains(str, ":")
}

func MatchIP6Addr(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

func MatchIPAddr(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}

func MatchPtr(str string) bool {
	return NsptrReg.MatchString(str)
}

func MatchName(str string) bool {
	return NameReg.MatchString(str)
}

func MatchDomainName(str string) bool {
	if str == "" || len(strings.ReplaceAll(str, ".", "")) > 255 {
		return false
	}
	return !MatchIPAddr(str) && DomainnameReg.MatchString(str)
}

func MatchDomainSRV(str string) bool {
	if !MatchDomainName(str) {
		return false
	}

	// Ref: https://tools.ietf.org/html/rfc2782
	//
	//	_Service._Proto.Name
	parts := strings.SplitN(str, ".", 3)
	if len(parts) != 3 {
		return false
	}
	for i := 0; i < 2; i++ {
		if len(parts[i]) < 2 || parts[i][0] != '_' {
			return false
		}
	}
	return len(parts[2]) != 0
}

func MatchSize(str string) bool {
	return SizeReg.MatchString(str)
}

func MatchMonth(str string) bool {
	return MonthReg.MatchString(str)
}

func MatchDate(str string) bool {
	return DateReg.MatchString(str)
}

func MatchDateCompact(str string) bool {
	return DateCompactReg.MatchString(str)
}

func MatchZStackTime(str string) bool {
	return ZstackTimeReg.MatchString(str)
}

func MatchISOTime(str string) bool {
	return IsoTimeReg.MatchString(str)
}

func MatchISONoSecondTime(str string) bool {
	return IsoNoSecondTimeReg.MatchString(str)
}

func MatchFullISOTime(str string) bool {
	return FullisoTimeReg.MatchString(str)
}

func MatchISOTime2(str string) bool {
	return IsoTimeReg2.MatchString(str)
}

func MatchISONoSecondTime2(str string) bool {
	return IsoNoSecondTimeReg2.MatchString(str)
}

func MatchFullISOTime2(str string) bool {
	return FullisoTimeReg2.MatchString(str)
}

func MatchCompactTime(str string) bool {
	return CompactTimeReg.MatchString(str)
}

func MatchMySQLTime(str string) bool {
	return MysqlTimeReg.MatchString(str)
}

func MatchNormalTime(str string) bool {
	return NormalTimeReg.MatchString(str)
}

func MatchFullNormalTime(str string) bool {
	return FullnormalTimeReg.MatchString(str)
}

func MatchRFC2882Time(str string) bool {
	return Rfc2882TimeReg.MatchString(str)
}

func MatchEmail(str string) bool {
	return EmailReg.MatchString(str)
}

func MatchMobile(str string) bool {
	return ChinaMobileReg.MatchString(str)
}

func MatchFS(str string) bool {
	return FsFormatReg.MatchString(str)
}

func MatchUSCurrency(str string) bool {
	return UsCurrencyReg.MatchString(str)
}

func MatchEUCurrency(str string) bool {
	return EuCurrencyReg.MatchString(str)
}
