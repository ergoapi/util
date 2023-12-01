package exnet

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/go-connections/nat"
)

// LocalIP 获取本机 ip
// 获取第一个非 loopback ip
func LocalIP() (net.IP, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				return v4, nil
			}
		}
	}
	return nil, errors.New("cannot find local IP address")
}

// Lists 获取本机非loopback ip
func Lists() (*[]net.Addr, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var allAddrs []net.Addr
	for i := 0; i < len(tables); i++ {
		if (tables[i].Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, _ := tables[i].Addrs()
		for j := 0; j < len(addrs); j++ {
			allAddrs = append(allAddrs, addrs[j])
		}
	}
	return &allAddrs, nil
}

// LocalIPs 获取本机非loopback ip
func LocalIPs() (addr []string) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				addr = append(addr, v4.String())
			}
		}
	}
	return addr
}

// GetMacAddrs 获取本机的Mac网卡地址列表.
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, _ := net.Interfaces()
	if len(netInterfaces) > 0 {
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			macAddrs = append(macAddrs, macAddr)
		}
	}
	return
}

func IsLocalHostAddrs() (*[]net.Addr, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var allAddrs []net.Addr
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := netInterfaces[i].Addrs()
		if err != nil {
			continue
		}
		for j := 0; j < len(addrs); j++ {
			allAddrs = append(allAddrs, addrs[j])
		}
	}
	return &allAddrs, nil
}

// GetFreePort 获取空闲端口
func GetFreePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return port
}

// IsPrivateIP 私网ip
func IsPrivateIP(str string) (bool, error) {
	ip := net.ParseIP(str)
	if ip == nil {
		return false, errors.New("str is not valid ip")
	}
	return IsPrivateNetIP(ip), nil
}

// IsPrivateNetIP 私网ip
func IsPrivateNetIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return true
}

// CheckHostPort if a port is available
func CheckHostPort(host string, port int) (status bool, err error) {
	// Concatenate a colon and the port
	host = host + ":" + strconv.Itoa(port)

	// Try to create a server with the port
	server, err := net.Listen("tcp", host)

	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}

	// close the server
	_ = server.Close()

	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil
}

// Check if a port is available
func Check(port int) (status bool, err error) {
	return CheckHostPort("", port)
}

// Deprecated: Use OutboundIPv2 instead.
// OutboundIP 获取本机的出口IP.
func OutboundIP() (string, error) {
	res := ""
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if conn != nil {
		addr := conn.LocalAddr().(*net.UDPAddr)
		res = addr.IP.String()
		_ = conn.Close()
	}
	return res, err
}

// OutboundIPv2 获取出口IP
func OutboundIPv2() (string, error) {
	resp, err := http.Get("https://ip.ysicing.cloud")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body := resp.Body
	txt, err := io.ReadAll(body)

	if err != nil {
		return "", err
	}

	ip := string(txt)

	if CheckIP(ip) == false {
		return ip, errors.New(ip + ": not ipv4")
	}

	return ip, nil
}

func GetAddrPort(addrPort string) (string, int) {
	parts := strings.Split(addrPort, ":")
	port, _ := strconv.Atoi(parts[1])
	return parts[0], port
}

func IsLocalIP(ip string, addrs *[]net.Addr) bool {
	for _, address := range *addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && ipnet.IP.String() == ip {
			return true
		}
	}
	return false
}

func CheckIP(i string) bool {
	if !strings.Contains(i, ":") {
		return net.ParseIP(i) != nil
	}
	return false
}

// GetIPByHostname 返回主机名对应的 IPv4地址.
func GetIPByHostname(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", nil
	}
	return "", err
}

// GetIpsByDomain 获取互联网域名/主机名对应的 IPv4 地址列表.
func GetIpsByDomain(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if ips != nil {
		var ipstrs []string
		for _, v := range ips {
			if v.To4() != nil {
				ipstrs = append(ipstrs, v.String())
			}
		}
		return ipstrs, nil
	}
	return nil, err
}

// GetHostByIP 获取指定的IP地址对应的主机名.
func GetHostByIP(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}

var equalHostIPs = map[string]interface{}{
	"":          nil,
	"127.0.0.1": nil,
	"0.0.0.0":   nil,
	"localhost": nil,
}

func IsPortBindingEqual(a, b nat.PortBinding) bool {
	if a.HostPort == b.HostPort {
		if _, ok := equalHostIPs[a.HostIP]; ok {
			if _, ok := equalHostIPs[b.HostIP]; ok {
				return true
			}
		}
	}
	return false
}
