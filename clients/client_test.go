/*
Tests for the clients package
*/
package clients

import (
	"fmt"
	"testing"
)

func TestExternalClientFromUrl(t *testing.T) {
	ip := "10.0.20.11"
	port := int64(8551)
	ext, err := ExternalClientFromURL(fmt.Sprintf("%s:%d", ip, port), "client")
	if err != nil {
		t.Fatal(err)
	}
	if ext.GetIP().String() != ip {
		t.Fatalf("Incorrect IP: want %s, got %s", ip, ext.IP.String())
	}
	extPort := ext.GetPort()
	if extPort == nil || *extPort != port {
		t.Fatalf("Incorrect port: want %d, got %d", port, extPort)
	}

	// Try without port
	ext, err = ExternalClientFromURL(ip, "client")
	if err != nil {
		t.Fatal(err)
	}
	if ext.GetIP().String() != ip {
		t.Fatalf("Incorrect IP: want %s, got %s", ip, ext.IP.String())
	}
	extPort = ext.GetPort()
	if extPort != nil {
		t.Fatalf("Incorrect port: want nil, got %d", extPort)
	}
}
