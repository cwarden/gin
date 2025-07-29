package gin

import (
	"fmt"
	"net"
)

// CheckPort checks if a port is available for binding
func CheckPort(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// FindAvailablePort finds an available port starting from the given port
func FindAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if CheckPort(port) {
			return port
		}
	}
	return -1
}

// FindAvailablePorts finds two consecutive available ports
func FindAvailablePorts(proxyPort, appPort int) (int, int) {
	// First try the requested ports if they're different
	if proxyPort != appPort && CheckPort(proxyPort) && CheckPort(appPort) {
		return proxyPort, appPort
	}

	// Find available proxy port
	availableProxyPort := proxyPort
	if !CheckPort(proxyPort) {
		availableProxyPort = FindAvailablePort(proxyPort)
		if availableProxyPort == -1 {
			return -1, -1
		}
	}

	// Find available app port (ensure it's different from proxy port)
	availableAppPort := appPort
	if !CheckPort(appPort) || appPort == availableProxyPort {
		// Start searching from appPort, but skip the proxy port
		for port := appPort; port < appPort+100; port++ {
			if port != availableProxyPort && CheckPort(port) {
				availableAppPort = port
				break
			}
		}
		if availableAppPort == appPort && (!CheckPort(appPort) || appPort == availableProxyPort) {
			return -1, -1
		}
	}

	return availableProxyPort, availableAppPort
}
