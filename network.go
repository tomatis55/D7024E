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
	count int
}

func Listen(ip string, port int) {
	// will listen for udp-packets on the provided ip and port
	// when a packet is detected start a goRoutine to handle it

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	sock, _ := net.ListenUDP("udp", &addr)
	defer sock.Close()

	for {
		buf := make([]byte, 1024)
		rlen, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		go handlePacket(buf, rlen)
	}
}

func handlePacket(buf []byte, rlen int) {
	// do different stuff depending on message
	// currently we only print the data in it
	fmt.Println(string(buf[0:rlen]))

	/*
		Ping
		FindContact
		FindData
		Store
	*/

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
