// package main

// import "time"

// func main() {
// 	id := NewKademliaID("5465747261687964726F63616E6E6162696E6F6C")
// 	// id := NewKademliaID("0000000000000000000000000000000000000000")
// 	d := id.CalcDistance(id)
// 	me := Contact{ID: id, Address: "127.0.0.1:2000", distance: d}
// 	kad := Kademlia{RoutingTable: NewRoutingTable(me), K: 3, Data: make(map[string][]byte)}

// 	n := Network{kad}

// 	go n.Listen("127.0.0.1", 2000)
// 	Pinger(n, me)
// }

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello world")
	arg := os.Args

	idSuperNode := "0000000000000000000000000000000000000000"
	ipSuperNode := "172.20.0.2"
	portSuperNode := ":81"
	portNode := ":82"
	ip := arg[1]

	network := &Network{}
	_ = &Contact{}

	// Initialize the super node
	if ip == ipSuperNode {
		me := NewContact(NewKademliaID(idSuperNode), ipSuperNode+portSuperNode)
		network = &Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}

		network.Listen(ip, 81) // run on c0/super node and run the code below to test ping in another container

		// Initialize the node and add the super node as a contact, then send a msg to let other nodes know of its existance
	} else {
		me := NewContact(NewRandomKademliaID(), ip+portNode)
		network = &Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}
		// add supernode as contact
		network.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode+portSuperNode))
		network.SendFindContactMessage(&me)

		contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode+portSuperNode)

		// test function to see if the super node is added as a contact
		retContact := network.Kademlia.LookupContact(&contact)
		if retContact[0].Address == contact.Address {
			fmt.Println("Jag existerar!")
		} else {
			fmt.Println(":(((((")
		}

		Pinger(*network, contact)
		// network.SendPingMessage(&contact)
		// network.SendPingMessage(&contact)
		// network.SendPingMessage(&contact)
		// network.SendPingMessage(&contact)

	}

}

func Pinger(n Network, me Contact) {

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		n.SendPingMessage(&me)
		time.Sleep(3 * time.Second)
	}

}
