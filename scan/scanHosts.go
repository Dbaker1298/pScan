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
