package scan_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Dbaker1298/pScan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			// Initialize the list with a host.
			if err := hl.Add("host1"); err != nil {
				t.Fatalf("failed to initialize list: %v", err)
			}

			err := hl.Add(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectErr)
				}

				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("expected error: %q, got: %q", tc.expectErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length: %d, got: %d instead\n", tc.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf("Expected host name: %q, as index 1, but got: %q instead\n", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			// Initialize the list with a host.
			for _, host := range []string{"host1", "host2"} {
				if err := hl.Add(host); err != nil {
					t.Fatalf("failed to initialize list: %v", err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil instead\n", tc.expectErr)
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error: %q, got: %q instead\n", tc.expectErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length: %d, got: %d instead\n", tc.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[0] == tc.host {
				t.Errorf("Host name %q should NOT be in the list, but it is\n", tc.host)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}

	hostName := "host1"
	hl1.Add(hostName)

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("failed to save hosts list: %s", err)
	}

	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("Host %q should match %q host.\n", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}

	if err := os.Remove(tf.Name()); err != nil {
		t.Errorf("Error deleting temp file: %s", err)

		hl := &scan.HostsList{}

		if err := hl.Load(tf.Name()); err != nil {
			t.Errorf("Expected no error, got: %q instead\n", err)
		}
	}
}
