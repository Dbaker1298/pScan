package scan_test

import (
	"errors"
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

			if hl.Hosts[0] != tc.host {
				t.Errorf("Host name %q should NOT be int list, but it is\n", tc.host)
			}
		})
	}
}
