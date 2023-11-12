package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Dbaker1298/pScan/scan"
)

// Since this app saves the hosts list to a file, we need to create a temporary file.
// Let's create an auxiliary function to set the test environment up.
func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create a temporary file.
	tf, err := ioutil.TempFile("", "pScan")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	tf.Close()

	// Initialize the list if requested.
	if initList {
		hl := &scan.HostsList{}

		for _, host := range hosts {
			hl.Add(host)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatalf("failed to save hosts list: %v", err)
		}

	}

	// Return the temporary file name and a function to clean up.
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	// Define hosts for actions test
	hosts := []string{"host1", "host2", "host3"}

	// Define test cases
	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up Action test
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			// Define var to capture Action output
			var out bytes.Buffer

			// Execute Action and capture output
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("Expected no error, got: %q\n", err)
			}

			// Compare output
			if out.String() != tc.expectedOut {
				t.Errorf("Expected output: %q, got: %q instead\n", tc.expectedOut, out.String())
			}
		})
	}
}

func TestScanAction(t *testing.T) {
	// Define hosts for scan action test
	hosts := []string{"localhost", "unknownhostoutthere"}

	// Setup can test
	tf, cleanup := setup(t, hosts, true)
	defer cleanup()

	ports := []int{}

	// Init port, 1 open, 1 closed
	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort(hosts[0], "0"))
		if err != nil {
			t.Fatalf("Failed to listen on port: %v\n", err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatalf("Failed to split host and port: %v\n", err)
		}

		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatalf("Failed to convert port string to int: %v\n", err)
		}

		ports = append(ports, portInt)

		if i == 1 {
			ln.Close()
		}
	}

	// Define expected output for scan action
	expectedout := fmt.Sprintln("localhost:")
	expectedout += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedout += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedout += fmt.Sprintln()
	expectedout += fmt.Sprintln("unknownhostoutthere: Host not found")
	expectedout += fmt.Sprintln()

	// Define var to capture Action output

	var out bytes.Buffer

	// Execute Action and capture output
	if err := scanAction(&out, tf, ports); err != nil {
		t.Fatalf("Expected no error, got: %q\n", err)
	}

	// Test scan output
	if out.String() != expectedout {
		t.Errorf("Expected output: %q, got: %q instead\n", expectedout, out.String())
	}
}

// Let's now add an integration test. The goal is to execute all commands
// in sequence, simulating a real user interaction. The user will add three
// hosts, list them, delete a host, and list them again.
func TestIntegration(t *testing.T) {
	// Define hosts for integration test
	hosts := []string{"host1", "host2", "host3"}

	// Set up integration test
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"

	hostsEnd := []string{"host1", "host3"}

	// Define var to capture Action output
	var out bytes.Buffer

	// Define expected output for all actions
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	// Execute all actions in the defined order; add -> list -> delete -> list
	// Add hosts to the list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no error, got: %q\n", err)
	}

	// List hosts
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got: %q\n", err)
	}

	// Delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Expected no error, got: %q\n", err)
	}

	// List hosts after deleting host2
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got: %q\n", err)
	}

	// Compare output
	if out.String() != expectedOut {
		t.Errorf("Expected output: %q, got: %q instead\n", expectedOut, out.String())
	}
}
