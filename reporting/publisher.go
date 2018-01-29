package reporting

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel      *amqp.Channel
	connexion    *amqp.Connection
	exchangeName string
	routingKey   string
}

func (pub *Publisher) Init(params map[string]string) {
	var err error

	pub.connexion, err = amqp.Dial(params["url"])
	failOnError(err, "Failed to connect to RabbitMQ")

	pub.channel, err = pub.connexion.Channel()
	failOnError(err, "Failed to open a channel")

	pub.exchangeName = params["exchange_name"]
	pub.routingKey = params["routing_key"]

	err = pub.channel.ExchangeDeclare(
		pub.exchangeName, // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

func (pub *Publisher) Close() {
	pub.connexion.Close()
	pub.channel.Close()
}

func (pub *Publisher) Send(report Report) {
	marshal, err := report.Encode()
	if err != nil {
		fmt.Println("Error encoding report", err.Error())
		return
	}
	pub.channel.Publish(
		pub.exchangeName, // exchange
		pub.routingKey,   // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(marshal),
		})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		time.Sleep(20 * time.Second)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
