package scan

import (
	"fmt"
	"net"
	"time"
)

// Define a new custom type PortState that represents the state for
// single TCP port.
type PortState struct {
	Port int
	Open state
}

// Define custom type `state` as a wrapper for the `bool` type.
type state bool

// String converts the boolean value of the state type to human-readable string.
// Using the Stringer interface to implement the String method.
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// scanPort perfoms a TCP scan on a single port.
func scanPort(host string, port int) PortState {
	p := PortState{Port: port}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	// Verify the function returned an error. If so, assume the port is closed.
	// This is a naive approach, but it works for our purposes.
	if err != nil {
		return p
	}

	// Close the connection if it was successful. Set the property to true.
	scanConn.Close()
	p.Open = true

	return p
}

// The scanPort function is private. We do not want users to call it directly.
// Instead, we will create a public function that will call scanPort for each
// port we want to scan.

// Results represents the results of a port scan for a single host.
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// Run perfoms a TCP scan on the hosts list
func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	for _, host := range hl.Hosts {
		r := Results{Host: host}

		// Use ne.LookupHost to resolve the hostname to an IP address.
		// If the host is not found, set the NotFound property to true.
		if _, err := net.LookupHost(host); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		// If host was found, loop through the ports and call scanPort for each port.
		for _, port := range ports {
			r.PortStates = append(r.PortStates, scanPort(host, port))
		}

		res = append(res, r)
	}

	return res
}
