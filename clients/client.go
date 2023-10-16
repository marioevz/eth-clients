/*
Generic client interface, used to describe an abstract client, whether it be
an execution client or a consensus client
*/
package clients

import (
	"net"
	"strconv"
)

type Client interface {
	IsRunning() bool
	GetHost() string
	GetIP() net.IP
	ClientType() string
}

type ManagedClient interface {
	Client
	AddStartOption(...interface{})
	Start() error
	Shutdown() error
}

var _ Client = &ExternalClient{}

type ExternalClient struct {
	Type     string
	Host     string
	IP       net.IP
	Port     *int64
	EnodeURL string
}

func ExternalClientFromURL(url string, typ string) (*ExternalClient, error) {
	host, portStr, err := net.SplitHostPort(url)
	if err != nil {
		if errP, ok := err.(*net.AddrError); ok {
			if errP.Err == "missing port in address" {
				host = url
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var port *int64
	if portStr != "" {
		portint, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			return nil, err
		}
		port = &portint
	}
	return &ExternalClient{
		Type: typ,
		Host: host,
		IP:   net.ParseIP(host),
		Port: port,
	}, nil
}

func (m *ExternalClient) IsRunning() bool {
	// We can try pinging a certain port for status
	return true
}

func (m *ExternalClient) GetHost() string {
	return m.Host
}

func (m *ExternalClient) GetIP() net.IP {
	return m.IP
}

func (m *ExternalClient) GetPort() *int64 {
	return m.Port
}

func (m *ExternalClient) ClientType() string {
	return m.Type
}

func (m *ExternalClient) GetEnodeURL() (string, error) {
	return m.EnodeURL, nil
}
