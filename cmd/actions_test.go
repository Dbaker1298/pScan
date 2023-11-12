package cmd

import (
	"io/ioutil"
	"os"
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
