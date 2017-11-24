package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/services"
	"gotrading/strategies"

	"github.com/streadway/amqp"
	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges/kraken"
	"github.com/thrasher-/gocryptotrader/exchanges/liqui"

	"flag"
)

var (
	uri          = flag.String("uri", "amqp://developer:xLae4pzT@hc-amqp.dev:5672/hc", "AMQP URI")
	exchangeName = flag.String("exchange", "arbitrage.fanout", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "fanout", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {

	cfg := &config.Cfg
	err := cfg.LoadConfig("config.dat")
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)

	liquiEngine := new(liqui.Liqui)
	krakenEngine := new(kraken.Kraken)
	// bittrexEngine := new(bittrex.Bittrex)
	// gdaxEngine := new(gdax.GDAX)
	// poloniexEngine := new(poloniex.Poloniex)

	liqui := services.LoadExchange(cfg, "Liqui", liquiEngine)
	kraken := services.LoadExchange(cfg, "Kraken", krakenEngine)
	// bittrex := services.LoadExchange(cfg, "Bittrex", bittrexEngine)
	// poloniex := services.LoadExchange(cfg, "Poloniex", poloniexEngine)
	// gdax := services.LoadExchange(cfg, "GDAX", gdaxEngine)

	// exchanges := []core.Exchange{kraken, liqui, gdax, bittrex}
	exchanges := []core.Exchange{kraken, liqui}

	mashup := core.ExchangeMashup{}
	mashup.Init(exchanges)

	from := core.Currency("BTC")
	to := from
	depth := 3
	nodes, paths, _ := graph.PathFinder(mashup, from, to, depth)

	// Create a map
	endpointLookup := make(map[string][]graph.EndpointLookup)
	for _, n := range nodes {
		paths := paths[n.ID()]
		lookups, ok := endpointLookup[n.Exchange.Name]
		if !ok {
			lookups = make([]graph.EndpointLookup, 0)
		}
		lookup := graph.EndpointLookup{n, len(paths)}
		endpointLookup[n.Exchange.Name] = append(lookups, lookup)
	}
	for _, exch := range exchanges {
		endpointLookup[exch.Name] = graph.MergeSort(endpointLookup[exch.Name])
	}

	arbitrage := strategies.Arbitrage{}

	delayBetweenReqs := make(map[string]time.Duration, len(exchanges))
	delayBetweenReqs["Kraken"] = time.Duration(500)
	delayBetweenReqs["Liqui"] = time.Duration(500)

	conn, err := amqp.Dial("amqp://developer:xLae4pzT@hc-amqp.dev:5672/hc")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"arbitrage.routing", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	for _, exch := range exchanges {
		nodes := endpointLookup[exch.Name]
		go services.StartPollingOrderbooks(exch, nodes, delayBetweenReqs[exch.Name], func(n graph.Endpoint) {
			chains := arbitrage.Run(paths[n.ID()])
			rows := make([][]string, 0)
			for _, chain := range chains {
				if chain.Performance == 0 {
					continue
				}
				ordersCount := len(chain.Path.Nodes)
				row := make([]string, ordersCount+5)
				for j, node := range chain.Path.Nodes {
					row[j] = node.Description()
				}

				row[ordersCount] = strconv.FormatFloat(chain.Performance, 'f', 6, 64)
				row[ordersCount+1] = strconv.FormatFloat(chain.Volume, 'f', 6, 64)
				row[ordersCount+2] = strconv.FormatFloat(chain.Volume*chain.Performance, 'f', 6, 64)
				row[ordersCount+3] = strconv.FormatFloat(chain.Volume*chain.Performance-chain.Volume, 'f', 6, 64)

				t := time.Now()
				row[ordersCount+4] = t.Format("2006-01-02 15:04:05")
				rows = append(rows, row)
				fmt.Println(strings.Join(row[:], ","))

				marshal, _ := json.Marshal(chain)

				err = ch.Publish(
					"arbitrage.routing", // exchange
					"usd.btc",           // routing key
					false,               // mandatory
					false,               // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(marshal),
					})
			}
			if len(rows) > 0 {
				// table.AppendBulk(rows)
				// table.Render()
			}
		})
	}

	<-interrupt
}

// This example declares a durable Exchange, and publishes a single message to
// that Exchange with a given routing key.
//

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
