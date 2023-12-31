package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/Dbaker1298/pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", ps.Open.String())
	}

	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	// Testing against localhost
	host := "localhost"
	hl := &scan.HostsList{}

	hl.Add(host)

	ports := []int{}

	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatalf("Failed to listen on port: %v\n", err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatalf("Failed to split host and port: %v\n", err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatalf("Failed to convert port string to int: %v\n", err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)

	// Verify the results of HostFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Errorf("Expected host %q to be found, but it was not\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[0], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectedState)
		}
	}
}

// Test when the host is not found
func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hl := &scan.HostsList{}

	hl.Add(host)

	res := scan.Run(hl, []int{})

	// Verify the results of HostNotFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	if !res[0].NotFound {
		t.Errorf("Expected host %q to NOT be found, but it was\n", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Errorf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
