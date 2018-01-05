package networking

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Gatling struct {
	Clients                              []*http.Client
	MaxRequestsPerSecondsForHost         map[string]int
	LastRequestFromClientToHostOccuredAt map[*http.Client]map[string]time.Time
	DefaultMaxRequestsPerSecondsForHost  int
	Mutexes                              map[*http.Client]sync.RWMutex
	IsVerbose                            bool
	RoundRobin                           int
}

var instance *Gatling
var once sync.Once

func SharedGatling() *Gatling {
	once.Do(func() {
		instance = &Gatling{}
		instance.warmUp()
	})
	return instance
}

func (g *Gatling) warmUp() {

	g.LastRequestFromClientToHostOccuredAt = make(map[*http.Client]map[string]time.Time)
	g.DefaultMaxRequestsPerSecondsForHost = 3
	g.Mutexes = make(map[*http.Client]sync.RWMutex)

	addrs, _ := net.InterfaceAddrs()
	eligibleAddrs := []net.Addr{}
	for _, addr := range addrs {
		if strings.HasPrefix(addr.String(), "10.0.") {
			eligibleAddrs = append(eligibleAddrs, addr)
		}
	}

	fmt.Println(eligibleAddrs)
	if len(eligibleAddrs) > 0 {
		g.Clients = make([]*http.Client, len(eligibleAddrs))
		for i, addr := range eligibleAddrs {

			tr := &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			}

			tcpAddr := &net.TCPAddr{
				IP: addr.(*net.IPNet).IP,
			}

			d := net.Dialer{LocalAddr: tcpAddr}

			tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				fmt.Println(d.LocalAddr, addr)
				return d.DialContext(ctx, network, addr)
			}
			g.Clients[i] = &http.Client{Transport: tr}
		}
	} else {
		g.Clients = make([]*http.Client, 1)

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}

		ief, err := net.InterfaceByName("en0")
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
		g.Clients[0] = &http.Client{Transport: tr}
	}
}

func (g *Gatling) GET(exchURL string) ([]byte, error) {

	if g.IsVerbose {
		log.Println("Gatling> Preparing interface", exchURL)
	}

	client := g.Clients[0]

	var contents []byte
	var err error

	mutex := g.Mutexes[client]
	URL, err := url.Parse(exchURL)
	hostname := URL.Hostname()
	maxRequestsPerSeconds := g.MaxRequestsPerSecondsForHost[hostname]
	if maxRequestsPerSeconds == 0 {
		maxRequestsPerSeconds = g.DefaultMaxRequestsPerSecondsForHost
	}
	delayBetweenRequests := 1.0 / float64(maxRequestsPerSeconds)
	mutex.RLock()
	lastOccurence, ok := g.LastRequestFromClientToHostOccuredAt[client][hostname]
	mutex.RUnlock()

	t := delayBetweenRequests - time.Since(lastOccurence).Seconds()
	if ok && t > 0 {
		time.Sleep(time.Duration(t*1000) * time.Millisecond)
	}

	mutex.Lock()
	if _, ok := g.LastRequestFromClientToHostOccuredAt[client]; ok {
		g.LastRequestFromClientToHostOccuredAt[client][hostname] = time.Now()

	} else {
		g.LastRequestFromClientToHostOccuredAt[client] = make(map[string]time.Time)
		g.LastRequestFromClientToHostOccuredAt[client][hostname] = time.Now()
	}
	mutex.Unlock()

	res, err := client.Get(exchURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("common.SendHTTPGetRequest() error: HTTP status code %d", res.StatusCode)
	}

	contents, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return contents, err
}