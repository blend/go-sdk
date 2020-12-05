package webutil

import (
	"net"
	"sync"
)

var (
	defaultGetInterfaceAddrs                            = net.InterfaceAddrs
	getInterfaceAddrs        func() ([]net.Addr, error) = defaultGetInterfaceAddrs
	getInterfaceAddrsLock    sync.Mutex
)

// LocalIP returns the local server ip.
func LocalIP() string {
	addrs, err := getInterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	// return the loopback ...
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// setGetInterfaceAddrs sets the `getInterfaceAddrs` factory.
func setGetInterfaceAddrs(gia func() ([]net.Addr, error)) {
	getInterfaceAddrsLock.Lock()
	defer getInterfaceAddrsLock.Unlock()
	getInterfaceAddrs = gia
}

// restoreGetInterfaceAddrs restores the `getInterfaceAddrs` factory.
func restoreGetInterfaceAddrs() {
	getInterfaceAddrsLock.Lock()
	defer getInterfaceAddrsLock.Unlock()
	getInterfaceAddrs = defaultGetInterfaceAddrs
}
