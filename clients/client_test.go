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
		url             string
		expectedAddress string
		expectedHost    string
		expectedIP      net.IP
		expectedPort    *int64
	}{
		{
			url:             "10.0.20.11:8551",
			expectedAddress: "http://10.0.20.11:8551",
			expectedHost:    "10.0.20.11",
			expectedIP:      net.ParseIP("10.0.20.11"),
			expectedPort:    pPort(8551),
		},
		{
			url:             "10.0.20.11",
			expectedAddress: "http://10.0.20.11",
			expectedHost:    "10.0.20.11",
			expectedIP:      net.ParseIP("10.0.20.11"),
			expectedPort:    nil,
		},
		{
			url:             "somehost:8551",
			expectedAddress: "http://somehost:8551",
			expectedHost:    "somehost",
			expectedIP:      nil,
			expectedPort:    pPort(8551),
		},
		{
			url:             "somehost",
			expectedAddress: "http://somehost",
			expectedHost:    "somehost",
			expectedIP:      nil,
			expectedPort:    nil,
		},
		{
			url:             "https://somehost",
			expectedAddress: "https://somehost",
			expectedHost:    "somehost",
			expectedIP:      nil,
			expectedPort:    nil,
		},
		{
			url:             "https://somehost:1234",
			expectedAddress: "https://somehost:1234",
			expectedHost:    "somehost",
			expectedIP:      nil,
			expectedPort:    pPort(1234),
		},
		{
			url:             "https://user:pass@bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io",
			expectedAddress: "https://user:pass@bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io",
			expectedHost:    "bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io",
			expectedIP:      nil,
			expectedPort:    nil,
		},
		{
			url:             "https://user:pass@bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io/",
			expectedAddress: "https://user:pass@bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io",
			expectedHost:    "bn.lighthouse-geth-1.srv.dencun-devnet-9.ethpandaops.io",
			expectedIP:      nil,
			expectedPort:    nil,
		},
	} {
		ext, err := ExternalClientFromURL(test.url, "client")
		if err != nil {
			t.Fatal(err)
		}
		if ext.GetAddress() != test.expectedAddress {
			t.Fatalf("Incorrect address: want %s, got %s", test.expectedAddress, ext.GetAddress())
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
