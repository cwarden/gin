package gin

import (
	"fmt"
	"net"
	"testing"
)

func TestCheckPort(t *testing.T) {
	// Test with a port that should be available
	port := 54321
	if !CheckPort(port) {
		t.Errorf("Expected port %d to be available, but it wasn't", port)
	}

	// Test with a port that we'll occupy
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	defer ln.Close()

	if CheckPort(port) {
		t.Errorf("Expected port %d to be unavailable, but it was reported as available", port)
	}
}

func TestFindAvailablePort(t *testing.T) {
	// Test finding an available port from a starting point
	startPort := 54300
	availablePort := FindAvailablePort(startPort)

	if availablePort == -1 {
		t.Error("Failed to find an available port")
	}

	if availablePort < startPort || availablePort >= startPort+100 {
		t.Errorf("Available port %d is outside expected range [%d, %d)", availablePort, startPort, startPort+100)
	}

	// Verify the found port is actually available
	if !CheckPort(availablePort) {
		t.Errorf("FindAvailablePort returned port %d, but it's not available", availablePort)
	}
}

func TestFindAvailablePortWithOccupiedPorts(t *testing.T) {
	// Occupy several ports
	startPort := 54400
	listeners := make([]net.Listener, 5)
	for i := 0; i < 5; i++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", startPort+i))
		if err != nil {
			t.Fatalf("Failed to listen on port %d: %v", startPort+i, err)
		}
		listeners[i] = ln
		defer ln.Close()
	}

	// Find should skip occupied ports
	availablePort := FindAvailablePort(startPort)
	if availablePort < startPort+5 {
		t.Errorf("Expected port to be at least %d, got %d", startPort+5, availablePort)
	}
}

func TestFindAvailablePorts(t *testing.T) {
	// Test with ports that should be available
	proxyPort := 54500
	appPort := 54501

	foundProxy, foundApp := FindAvailablePorts(proxyPort, appPort)

	if foundProxy == -1 || foundApp == -1 {
		t.Error("Failed to find available ports")
	}

	if foundProxy == foundApp {
		t.Error("Found the same port for proxy and app")
	}

	// When ports are available, they should be returned as-is
	if CheckPort(proxyPort) && CheckPort(appPort) {
		if foundProxy != proxyPort || foundApp != appPort {
			t.Errorf("Expected to get requested ports back when available, got %d,%d instead of %d,%d",
				foundProxy, foundApp, proxyPort, appPort)
		}
	}
}

func TestFindAvailablePortsWithOccupiedProxyPort(t *testing.T) {
	proxyPort := 54600
	appPort := 54601

	// Occupy the proxy port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		t.Fatalf("Failed to listen on port %d: %v", proxyPort, err)
	}
	defer ln.Close()

	foundProxy, foundApp := FindAvailablePorts(proxyPort, appPort)

	if foundProxy == -1 || foundApp == -1 {
		t.Error("Failed to find available ports")
	}

	if foundProxy == proxyPort {
		t.Error("Got occupied proxy port back")
	}

	if foundProxy == foundApp {
		t.Error("Found the same port for proxy and app")
	}
}

func TestFindAvailablePortsWithBothOccupied(t *testing.T) {
	proxyPort := 54700
	appPort := 54701

	// Occupy both ports
	ln1, err := net.Listen("tcp", fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		t.Fatalf("Failed to listen on port %d: %v", proxyPort, err)
	}
	defer ln1.Close()

	ln2, err := net.Listen("tcp", fmt.Sprintf(":%d", appPort))
	if err != nil {
		t.Fatalf("Failed to listen on port %d: %v", appPort, err)
	}
	defer ln2.Close()

	foundProxy, foundApp := FindAvailablePorts(proxyPort, appPort)

	if foundProxy == -1 || foundApp == -1 {
		t.Error("Failed to find available ports")
	}

	if foundProxy == proxyPort || foundProxy == appPort {
		t.Error("Got occupied port back for proxy")
	}

	if foundApp == proxyPort || foundApp == appPort {
		t.Error("Got occupied port back for app")
	}

	if foundProxy == foundApp {
		t.Error("Found the same port for proxy and app")
	}
}

func TestFindAvailablePortsEnsuresDifferentPorts(t *testing.T) {
	// Test case where requested ports are the same
	port := 54800

	foundProxy, foundApp := FindAvailablePorts(port, port)

	if foundProxy == -1 || foundApp == -1 {
		t.Error("Failed to find available ports")
	}

	if foundProxy == foundApp {
		t.Error("Found the same port for proxy and app when requesting same port")
	}
}
