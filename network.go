package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

/*
M1. Network formation. [5p]. Your nodes must be able to form networks as described in the Kademlia
	paper. Kademlia is a protocol for facilitating Distributed Hash Tables (DHTs). Concretely,
	the following aspects of the algorithm must be implemented:
	(a) Pinging. This means that you must implement and use the PING message.
	(b) Network joining. Given the IP address, and any other data you decide, of any single
		node, a node must be able to join or form a network with that node.
	(c) Node lookup. When part of a network, each node must be able to retrieve the contact
		information of any other node in the same network.
*/

type Network struct {
	Kademlia Kademlia
}

/*
	PING
	FIND_CONTACT
	FIND_DATA
	STORE
*/

type Message struct {
	RPCtype string
	Sender  Contact
	Message []byte
}

// will listen for udp-packets on the provided ip and port
// when a packet is detected start a goRoutine to handle it
func Listen(ip string, port int) {
	// set up our connection to listen for udp on the specified address
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, _ := net.ListenUDP("udp", &addr)
	dec := gob.NewDecoder(conn)

	// listen to our connection, relaying all recieved packets to a handlePacket() go routine for proper handling
	for {
		var msg Message
		err := dec.Decode(&msg)
		if err != nil {
			fmt.Println("Decode error:", err)
		}

		fmt.Println("decoded message: ", msg.Message)

		go handlePacket(msg)
	}
}

// handles a packet, doing what needs to be done and sending the correct messages depending on the type of message recieved
func handlePacket(msg Message) {

	switch msg.RPCtype {
	case "PING":
		fmt.Println("you can ping, you can jive, having the time of your life")
	case "FIND_CONTACT":
		fmt.Println("find me, find me, find me a contact after midnight")
	case "FIND_DATA":
		fmt.Println("data, data, data, must be funny in the rich mans world")
	case "STORE":
		fmt.Println("the winner stores it all, the loser has to fall")
	default:
		fmt.Println("oh no something unexpected happened")
	}

}

// sends a message (encoded as a []byte) over udp towards a target address
func (network *Network) sendMessage(addr string, msg Message) {
	conn, _ := net.Dial("udp", addr)

	enc := gob.NewEncoder(conn)
	err := enc.Encode(msg)
	if err != nil {
		fmt.Println("Encode error:", err)
	}

}

// you guessed it, this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) {
	msg := Message{
		Message: []byte("PING pong! this is a PING message!"),
		RPCtype: "PING",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	fmt.Println("original message:", msg.Message)
	network.sendMessage(contact.Address, msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// message := []byte("greetings traveler! this is a FIND_CONTACT message!")
	// network.sendMessage("FIND_CONTACT", contact.Address, message)
}

func (network *Network) SendFindDataMessage(hash string) {
	// message := []byte("well met friend! this is a FIND_DATA message!")
	// network.sendMessage("FIND_DATA", contact.Address, message)
}

func (network *Network) SendStoreMessage(data []byte) {
	// message := []byte("hello! this is a STORE message!")
	// network.sendMessage("STORE", contact.Address, message)
}
