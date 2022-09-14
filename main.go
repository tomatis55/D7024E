package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello world")
	arg := os.Args

	id := arg[1]
	ip := arg[2]

	me := NewContact(NewKademliaID(id), ip) // use random id if not first node
	network := &Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}

	if ip != "172.20.0.2" {
		// add supernode as contact
		network.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.2"))
		network.SendFindContactMessage(&me)
	}

	// Listen(ip, 2000)
	network.SendPingMessage(&me)

}

// To join the network, a node u must have a contact to an already participating
// node w. u inserts w into the appropriate k-bucket. u then performs a node lookup
// for its own node ID.Finally, u refreshes all k-buckets further away than its closest
// neighbor.During the refreshes, u both populates its own k-buckets and inserts
// itself into other nodes’ k-buckets as necessary
