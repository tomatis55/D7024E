package main

import (
	. "d7024e"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello world")
	arg := os.Args

	idSuperNode := "0000000000000000000000000000000000000000"
	ipSuperNode := "172.20.0.2"
	port := ":80"
	ip := arg[1]

	fmt.Println("my ip is: ", ip)

	network := &Network{}
	_ = &Contact{}

	// Initialize the super node
	if ip == ipSuperNode {
		me := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
		network = &Network{Kademlia: Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}, Alpha: 3, Channel: make(chan Message)}

		network.Listen(ipSuperNode, 80) // run on c0/super node and run the code below to test ping in another container

		// Initialize the node and add the super node as a contact, then send a msg to let other nodes know of its existance
	} else {
		me := NewContact(NewRandomKademliaID(), ip+port)
		network = &Network{Kademlia: Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}, Alpha: 3, Channel: make(chan Message)}
		// add supernode as contact
		network.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode+port))
		network.SendFindContactMessage(&me)

		contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)

		// test function to see if the super node is added as a contact
		retContact := network.Kademlia.LookupContact(&contact)
		if retContact.GetContacts(6)[0].Address == contact.Address {
			fmt.Println("Jag existerar!")
		} else {
			fmt.Println(":(((((")
		}

		go network.Listen(ip, 80)
		Pinger(*network, contact)
	}

}

func Pinger(n Network, me Contact) {

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		n.SendPingMessage(&me)
		time.Sleep(3 * time.Second)
	}

}
