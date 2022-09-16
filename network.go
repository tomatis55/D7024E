package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
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

type Message struct {
	RPCtype string // PING, PING_ACK, FIND_CONTACT, FIND_DATA, STORE
	Sender  Contact
	Message []byte
	// more?
}

/*












 */

// will listen for udp-packets on the provided ip and port
// when a packet is detected start a goRoutine to handle it
func (network *Network) Listen(ip string, port int) {
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
			fmt.Println("Read error:", err)
		}

		var msg Message
		json.Unmarshal(buf[:rlen], &msg)

		network.handlePacket(msg)
	}
}

/*












 */

func (network *Network) updateBucket(sender Contact) {

	bucket := *network.Kademlia.RoutingTable.buckets[network.Kademlia.RoutingTable.getBucketIndex(sender.ID)]

	if bucket.Len() <= bucketSize {
		network.Kademlia.RoutingTable.AddContact(sender)
	} else {
		// might still be in bucket
		closestContact := network.Kademlia.RoutingTable.FindClosestContacts(sender.ID, 1)
		id := NewKademliaID("0000000000000000000000000000000000000000")

		// find closest contact has 0 distance means we are already in the bucket
		if closestContact[0].distance.Equals(id) {
			network.Kademlia.RoutingTable.AddContact(sender)
		} else {
			// ping buckets head to see if alive
			response, _ := network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			if response != nil {
				return
			} else {
				network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
				network.Kademlia.RoutingTable.AddContact(sender)
			}
		}
	}
}

/*












 */

func (network *Network) handlePacket(msg Message) {

	switch msg.RPCtype {
	case "PING":
		fmt.Println("you can ping, you can jive, having the time of your life")

		// add sender to my bucket
		network.updateBucket(msg.Sender)

		// send ack back
		pingAck := Message{
			Message: []byte("ping PONG, i hear you!"),
			RPCtype: "PING_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
		}
		network.sendMessage(msg.Sender.Address, pingAck)

	case "PING_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)

		fmt.Println("ping ackked")

	case "FIND_CONTACT":
		// do kademlia stuff to find the contact
		// add contacts to buckets?
		fmt.Println("find me, find me, find me a contact after midnight")
	case "FIND_DATA":
		fmt.Println("data, data, data, must be funny in the rich mans world")
	case "STORE":
		fmt.Println("the winner stores it all, the loser has to fall")
	default:
		fmt.Println("oh no something unexpected happened")
	}

}

/*












 */
// sends a message and returns its message if any...
func (network *Network) sendMessage(addr string, msg Message) ([]byte, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("DIAL error:", err)
	}
	conn.SetDeadline(time.Now().Add(time.Second))
	defer conn.Close()

	marshalled_msg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Marshal error:", err)
	}
	_, err = conn.Write(marshalled_msg)
	if err != nil {
		fmt.Println("Write error:", err)
	}

	buff := make([]byte, 1024)
	len, err := conn.Read(buff)
	if err == nil {
		return buff[:len], nil
	} else {
		return nil, err
	}
}

/*












 */
// you guessed it, this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) ([]byte, error) {
	msg := Message{
		Message: []byte("PING pong! this is a PING message!"),
		RPCtype: "PING",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	return network.sendMessage(contact.Address, msg)
}

/*












 */

func (network *Network) SendFindContactMessage(contact *Contact) {
	// message := []byte("greetings traveler! this is a FIND_CONTACT message!")
	// network.sendMessage("FIND_CONTACT", contact.Address, message)
}

/*












 */

func (network *Network) SendFindDataMessage(hash string) {
	// message := []byte("well met friend! this is a FIND_DATA message!")
	// network.sendMessage("FIND_DATA", contact.Address, message)
}

/*












 */

func (network *Network) SendStoreMessage(data []byte) {
	// ok SendStoreMessage(data)
	// vars skriver jag?
	// vem blir skickad meddelandet?
}
