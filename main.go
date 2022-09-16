package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello world")
	arg := os.Args

	idSuperNode := "0000000000000000000000000000000000000000"
	ipSuperNode := "172.20.0.2"
	ip := arg[1]

	network := &Network{}
	_ = &Contact{}

	// Initialize the super node
	if ip == ipSuperNode {
		me := NewContact(NewKademliaID(idSuperNode), ipSuperNode)
		network = &Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}

		// Initialize the node and add the super node as a contact, then send a msg to let other nodes know of its existance
	} else {
		me := NewContact(NewRandomKademliaID(), ip)
		network = &Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}
		// add supernode as contact
		network.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode))
		network.SendFindContactMessage(&me)
	}

	Listen(ip, 2000) // run on c0/super node and run the code below to test ping in another container

	// contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode)

	// // test function to see if the super node is added as a contact
	// retContact := network.Kademlia.LookupContact(&contact)
	// if retContact[0].Address == contact.Address {
	// 	fmt.Println("Jag existerar!")
	// } else {
	// 	fmt.Println(":(((((")
	// }

	// network.SendPingMessage(&contact)
	// network.SendPingMessage(&contact)
	// network.SendPingMessage(&contact)
	// network.SendPingMessage(&contact)

}
