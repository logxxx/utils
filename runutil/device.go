package runutil

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"sort"
	"strings"
)

var (
	ip      string
	macAddr string
)

func init() {
	ips := GetLocalIPAddrs()
	if len(ips) > 0 {
		ip = ips[0]
	}

	macAddrs := GetMacAddrs()
	if len(macAddrs) > 0 {
		macAddr = macAddrs[0]
	}

}

type interfaceSlice []net.Interface

func (x interfaceSlice) Len() int { return len(x) }
func (x interfaceSlice) get(i int) int {
	names := []string{"eth", "wlan", "e", "w"}
	for j, _ := range names {
		if strings.HasPrefix(x[i].Name, names[j]) {
			return j
		}
	}
	return len(names)
}
func (x interfaceSlice) Less(i, j int) bool {
	return x.get(i) < x.get(j)

}

func (x interfaceSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// GetLocalIPAddrs 获取本机IP地址，排除回环地址
func GetLocalIPAddrs() (ips []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(fmt.Sprintf("GetLocalIPAddrs panic:%v", err))
		return ips
	}
	sort.Sort(interfaceSlice(netInterfaces))

	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func GetMacAddrs() (macAddrs []string) {
	envMacAddrs := os.Getenv("MAC_ADDRS")
	if envMacAddrs != "" {
		return strings.Split(envMacAddrs, ",")
	}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Errorf("GetMacAddrs net.Interfaces err:%v", err)
		return macAddrs
	}
	sort.Sort(interfaceSlice(netInterfaces))

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	log.Debugf("GetMacAddrs result:%v", macAddrs)
	return macAddrs
}

func GetIP() string {
	return ip
}

func GetMacAddr() string {
	return macAddr
}
