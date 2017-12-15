package gatling

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	clients []*http.Client

	// Transport Queue
	// LastRequestPerDomain map[string][]Path
}

func (c *Client) Init(ifaces []string) {
	c.clients = make([]*http.Client, len(ifaces))
	for i, iface := range ifaces {

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}

		ief, err := net.InterfaceByName(iface)
		if err != nil {
			log.Fatal(err)
		}

		addrs, err := ief.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		tcpAddr := &net.TCPAddr{
			IP: addrs[1].(*net.IPNet).IP,
		}

		d := net.Dialer{LocalAddr: tcpAddr}

		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.DialContext(ctx, network, addr)
		}
		c.clients[i] = &http.Client{Transport: tr}
	}
}

func (c *Client) WarmUp() {

}

func (c *Client) Fire() {
	resp, err := c.clients[0].Get("https://google.com")
	fmt.Println(resp, err)
}

func (c *Client) Reload() {

}


		common.SendHTTPGetRequest(req, true, l.Verbose, &response.Data)

