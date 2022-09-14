package main

import (
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
	kademlia Kademlia
}

type message struct {
	RPCtype string
	sender  Contact
	message []byte
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
	defer conn.Close()

	// listen to our connection, relaying all recieved packets to a handlePacket() go routine for proper handling
	buf := make([]byte, 1024)
	for {
		rlen, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		go handlePacket(buf, rlen)
	}
}

// handles a packet, doing what needs to be done and sending the correct messages depending on the type of message recieved
/*
	Ping
	FindContact
	FindData
	Store
*/
func handlePacket(buf []byte, rlen int) {
	// do different stuff depending on message
	// currently we only print the data in it
	fmt.Println(string(buf[0:rlen]))

}

// sends any message over udp towards a target address
// TODO!: NOT DONE YET
// + implement some kind of message struct with senderInfo, message and RPCtype
func (network *Network) sendMessage(addr string, msg message) {
	conn, _ := net.Dial("udp", addr)
	_, err := conn.Write(msg.message)
	if err != nil {
		fmt.Println(err)
	}
}

// you guessed it, this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) {
	msg := message{message: []byte("PING pong! this is a PING message!"), sender: network.kademlia.routingTable.me, RPCtype: "PING"}
	go network.sendMessage(contact.Address, msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// message := []byte("greetings traveler! this is a FIND_CONTACT message!")
	// go network.sendMessage("FIND_CONTACT", contact.Address, message)
}

func (network *Network) SendFindDataMessage(hash string) {
	// message := []byte("well met friend! this is a FIND_DATA message!")
	// go network.sendMessage("FIND_DATA", contact.Address, message)
}

func (network *Network) SendStoreMessage(data []byte) {
	// message := []byte("hello! this is a STORE message!")
	// go network.sendMessage("STORE", contact.Address, message)
}
