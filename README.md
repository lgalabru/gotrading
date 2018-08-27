
![Gotrading logo](./assets/gotrading.png)



## Getting Started

I started this side project on my spare time a few months ago. The idea was too build a library that would help getting started with algotrading, and give me the ability to pass orders on any exchange.

#### Features

- [x] Sending orders to Liqui
- [x] Sending orders to Binance
- [x] Working around API rate limits 😎

I spent some time writing and re-writing the abstraction for adding more exchanges, it should now take 3 to 4 hours (testing not included) to add a new exchange.

#### Status

I need a few days/weeks to clean the project, write some documentation, consolidate and write some tests, before being able to have this bot running 24/7.
I recently shifted my spare time on playing with machine learning, I'll probably go back to this project at some point (depending on how is the community answering), but can't give an ETA at this point.


## Demo

Short video of the UI I've been developing for visualizing and debugging my trades (built with React + D3, listening to RabbitMQ events sent by the bot).

[![Gotrading demo](https://img.youtube.com/vi/P-G78LB2LfU/0.jpg)](http://www.youtube.com/watch?v=P-G78LB2LfU)

This dashboard is a separate project that I can also open source, just open an issue on this repo if you're interested.


## Architecture

#### Worker
First, we need to instanciate the EC2 or `gatling`, that will be responsible for fetching the quotes and sending orders (ideally one per exchange). This node needs to be in the same / nearest datacenter hosting the exchange (`ap-northeast-1a` for Binance), in order to limit latency.

The configuration of this node is a bit special: we will be attaching as many Virtual Network Interfaces as possible (2 for a t2-micro), and attaching as many Elastic IP as possible (2 EIP / VNI for t2-micro). 

This part of the project is managed with Terraform (also in a separate repo, open a new issue if you're interested):

```
module "binance" {
  source = "../modules/instances/worker"

  instance_type = "t2.micro"
  ami = "ami-12572374"

  network_interfaces_count = 2
  ips_per_network_interface = 2

  availability_zone = "ap-northeast-1a"
}

```

Thanks to this tuning, instead of having our worker being limited to **5** req/sec, we can have **20** req/sec, and we could theoritically scale this limit to **3,750** req/sec with a `c5d.18xlarge`.

This parameter is important, since the first strategy implemented is an arbitrage using 3 quotes on one exchange (ex: BTC→ETH, ETH→TRX, TRX→BTC).
By fetching the 3 quotes 6 times per seconds, we are more reactive than the users getting the quotes using the websocket API.


## How to use

This project is designed in layers free of inter-dependencies, and you should be able to be used as a library. 

- core: types required for modelizing the space (Currency, Pair, Orderbook, etc) 
- networking: mainly abstracting the concept of `gatling`
- graph: utils for manipulating the graph of pairs available on one exchange.
- exchanges: where new exchanges should be added

You can also definitely have a look at the code available in `./strategies`, to have a better of how everything can be orchestrated.


## License

Author: Ludovic Galabru

This project is licensed under the MIT License.

## Credits

[@egonelbre](https://github.com/egonelbre/gophers) for the gopher :)

