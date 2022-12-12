# Simple p2p network 

This is a simple p2p network that uses a blockchain to store data. It is a simple implementation of a blockchain and is not intended to be used in production.
I use this project to learn about blockchain and p2p networks.

## Consensus
In this project, I use snowball consensus algorithm to reach consensus. Snowball is a consensus algorithm that is used in the avalanche protocol. It is a probabilistic consensus algorithm that is used to reach consensus in a p2p network. It is a simple algorithm that is easy to implement and understand. It is also a good algorithm to use to learn about consensus algorithms.

## Network
In the p2p network, each node will choose some nodes as their neighbors
```text
             +------+
    +------->+ node +<------+
    |        +------+       |
    |                       |
    v                       v
+---+--+                 +--+---+
| node |                 | node |
+---+--+                 +--+---+
    ^                       ^
    |                       |
    |        +------+       |
    +------->+ node +<------+
             +------+
```

## Peer
A peer is a node in the network. Each peer has a blockchain and a list of neighbors. The blockchain is used to store data and the neighbors are used to communicate with other peers.
We use gRPC to communicate between peers. Each peer has a gRPC server that listens for requests from other peers. Each peer also has a gRPC client that is used to send requests to other peers.

```text
                                gRPC
                                  |
                                  v       +----+
                       +-----connection-->+peer|
                       |                  +----+
                       |
                       |
                       |
                   +---+--+               +----+
+----Service---+-->+ node +--connection-->+peer|
|PeerManager   |   +---+--+               +----+
|MessageManager|       |
+--------------+       |
                       |
                       |                  +----+
                       +-----connection-->+peer|
                                          +----+
```

## Build and run
To build and run the project, you need to have go installed on your machine. You can download go from [here](https://golang.org/dl/).

To build the project, run the following command in the root directory of the project
```bash
go build -o build/startnode cmd/startnode/main.go
```

To run the project, run the following command in the root directory of the project
```bash
./build/startnode -port 5000 -neighbors localhost:5001,localhost:5002
```