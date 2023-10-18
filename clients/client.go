/*
Generic client interface, used to describe an abstract client, whether it be
an execution client or a consensus client
*/
package clients

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client interface {
	IsRunning() bool
	GetAddress() string
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
	Address  string
	Host     string
	IP       net.IP
	Port     *int64
	EnodeURL string
}

func ExternalClientFromURL(url string, typ string) (*ExternalClient, error) {
	address := url
	hostPortAuth := address
	{
		splitArr := strings.Split(hostPortAuth, "://")
		if len(splitArr) == 2 {
			hostPortAuth = splitArr[1]
		} else {
			address = fmt.Sprintf("http://%s", address)
		}
	}
	hostPort := hostPortAuth
	{
		splitArr := strings.Split(hostPortAuth, "@")
		if len(splitArr) == 2 {
			hostPort = splitArr[1]
		}
	}
	host, portStr, err := net.SplitHostPort(hostPort)
	if err != nil {
		if errP, ok := err.(*net.AddrError); ok {
			if errP.Err == "missing port in address" {
				host = hostPort
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
		Address: address,
		Type:    typ,
		Host:    host,
		IP:      net.ParseIP(host),
		Port:    port,
	}, nil
}

func (m *ExternalClient) IsRunning() bool {
	// We can try pinging a certain port for status
	return true
}

func (m *ExternalClient) GetAddress() string {
	if m.Address != "" {
		return m.Address
	}
	if m.Port != nil {
		return fmt.Sprintf(
			"http://%s:%d",
			m.GetHost(),
			m.GetPort(),
		)
	} else {
		return fmt.Sprintf(
			"http://%s",
			m.GetHost(),
		)
	}
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
