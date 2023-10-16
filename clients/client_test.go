/*
Tests for the clients package
*/
package clients

import (
	"net"
	"testing"
)

func TestExternalClientFromUrl(t *testing.T) {
	pPort := func(p int64) *int64 {
		return &p
	}

	for _, test := range []struct {
		url          string
		expectedHost string
		expectedIP   net.IP
		expectedPort *int64
	}{
		{
			url:          "10.0.20.11:8551",
			expectedHost: "10.0.20.11",
			expectedIP:   net.ParseIP("10.0.20.11"),
			expectedPort: pPort(8551),
		},
		{
			url:          "10.0.20.11",
			expectedHost: "10.0.20.11",
			expectedIP:   net.ParseIP("10.0.20.11"),
			expectedPort: nil,
		},
		{
			url:          "somehost:8551",
			expectedHost: "somehost",
			expectedIP:   nil,
			expectedPort: pPort(8551),
		},
		{
			url:          "somehost",
			expectedHost: "somehost",
			expectedIP:   nil,
			expectedPort: nil,
		},
	} {
		ext, err := ExternalClientFromURL(test.url, "client")
		if err != nil {
			t.Fatal(err)
		}
		if test.expectedIP != nil && ext.GetIP() != nil {
			if !test.expectedIP.Equal(ext.GetIP()) {
				t.Fatalf("Incorrect IP: want %s, got %s", test.expectedIP, ext.IP.String())
			}
		} else if test.expectedIP != nil && ext.GetIP() == nil {
			t.Fatalf("Incorrect IP: want %s, got %s", test.expectedIP, ext.IP.String())
		} else if test.expectedIP == nil && ext.GetIP() != nil {
			t.Fatalf("Incorrect IP: want %s, got %s", test.expectedIP, ext.IP.String())
		}
		if ext.GetHost() != test.expectedHost {
			t.Fatalf("Incorrect Host: want %s, got %s", test.expectedHost, ext.GetHost())
		}
		extPort := ext.GetPort()
		if extPort != nil && test.expectedPort != nil {
			if *extPort != *test.expectedPort {
				t.Fatalf("Incorrect port: want %d, got %d", *test.expectedPort, *extPort)
			}
		} else if extPort != nil && test.expectedPort == nil {
			t.Fatalf("Incorrect port: want nil, got %d", *extPort)
		} else if extPort == nil && test.expectedPort != nil {
			t.Fatalf("Incorrect port: want %d, got nil", *test.expectedPort)
		}
	}
}
