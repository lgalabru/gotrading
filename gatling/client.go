package gatling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gotrading/core"
	"gotrading/graph"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/thrasher-/gocryptotrader/currency/pair"
)

type Gatling struct {
	Clients                              []*http.Client
	MaxRequestsPerSecondsForHost         map[string]int
	LastRequestFromClientToHostOccuredAt map[*http.Client]map[string]time.Time
	DefaultMaxRequestsPerSecondsForHost  int
	Mutexes                              map[*http.Client]sync.RWMutex
}

func (g *Gatling) WarmUp() {

	g.LastRequestFromClientToHostOccuredAt = make(map[*http.Client]map[string]time.Time)
	g.DefaultMaxRequestsPerSecondsForHost = 5
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

type pathFetched func(path graph.Path)

type indexedNode struct {
	Index int
	Node  *graph.Node
}

func (g *Gatling) fireRequest(vertice *graph.Vertice, i int, c chan indexedNode) {
	n := vertice.Content
	cp := pair.NewCurrencyPair(string(n.Endpoint.From), string(n.Endpoint.To))
	exch := n.Endpoint.Exchange

	client := g.Clients[i%len(g.Clients)]

	type Orderbook struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
	}
	type Response struct {
		Data map[string]Orderbook
	}

	response := Response{}
	curr := fmt.Sprintf("%s", cp.Display("_", false))

	req := fmt.Sprintf("%s/%s/%s/%s?limit=5", "https://api.Liqui.io/api", "3", "depth", curr)

	t1 := time.Now()
	err := g.SendHTTPGetRequest(client, req, true, false, &response.Data)
	t2 := time.Now()
	src := response.Data[curr]

	if err == nil {
		dst := &core.Orderbook{}
		dst.CurrencyPair = core.CurrencyPair{n.Endpoint.From, n.Endpoint.To}
		dst.Bids = make([]core.Order, 0)
		dst.Asks = make([]core.Order, 0)
		dst.StartedLastUpdateAt = t1
		dst.EndedLastUpdateAt = t2

		for _, ask := range src.Asks {
			if exch.IsCurrencyPairNormalized == true {
				dst.Asks = append(dst.Asks, core.NewAsk(dst.CurrencyPair, ask[0], ask[1]))
			} else {
			}
		}
		for _, bid := range src.Bids {
			if exch.IsCurrencyPairNormalized == true {
				dst.Bids = append(dst.Bids, core.NewBid(dst.CurrencyPair, bid[0], bid[1]))
			} else {
			}
		}
		n.Endpoint.Orderbook = dst
	} else {
		fmt.Println("Error", n.Endpoint.Description(), err)
	}
	c <- indexedNode{i, n}
}

func (g *Gatling) FireRequests(vertices []*graph.Vertice, fn pathFetched) {
	path := graph.Path{}
	path.Nodes = make([]*graph.Node, len(vertices))
	c := make(chan indexedNode, len(vertices))

	fmt.Println("-----------------")

	for i, v := range vertices {
		if len(g.Clients) > 1 {
			go g.fireRequest(v, i, c)
		} else {
			g.fireRequest(v, i, c)
		}
	}
	for range vertices {
		indexedNode := <-c
		path.Nodes[indexedNode.Index] = indexedNode.Node
	}

	path.Encode()
	fn(path)
}

func (g *Gatling) SendHTTPGetRequest(client *http.Client, exchURL string, jsonDecode, isVerbose bool, result interface{}) error {

	if isVerbose {
		log.Println("Gatling> Preparing interface", exchURL)
	}

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

	if isVerbose {
		log.Println("Gatling> Fetching")
	}

	res, err := client.Get(exchURL)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("common.SendHTTPGetRequest() error: HTTP status code %d", res.StatusCode)
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// if isVerbose {
	// 	log.Println("Raw Resp: ", string(contents[:]))
	// }

	defer res.Body.Close()

	if jsonDecode {
		if !strings.Contains(reflect.ValueOf(result).Type().String(), "*") {
			return errors.New("json decode error - memory address not supplied")
		}

		err := json.Unmarshal(contents, result)
		if err != nil {
			log.Println(string(contents[:]))
			return err
		}
	}

	return nil
}

// common.SendHTTPGetRequest(req, true, l.Verbose, &response.Data)
