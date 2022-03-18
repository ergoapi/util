package exnet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

//LocalIP 获取本机 ip
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
	return nil, fmt.Errorf("cannot find local IP address")
}

//LocalIPs 获取本机非loopback ip
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

// GetFreePort 获取空闲端口
func GetFreePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		return 0
	}

	return port
}

// IsPrivateIP 私网ip
func IsPrivateIP(ip net.IP) bool {
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
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
