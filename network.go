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
	alpha    int
}

type Message struct {
	RPCtype      string // PING, PING_ACK, FIND_CONTACT, FIND_CONTACT_ACK, FIND_DATA, FIND_DATA_ACK, STORE, STORE_ACK
	Sender       Contact
	Message      []byte
	QueryContact *Contact
	Hash         string
	Data         []byte
	// more?
}

// will listen for udp-packets on the provided ip and port
// when a packet is detected start a goRoutine to handle it
func (network *Network) Listen(ip string, port int) {
	// set up our connection to listen for udp on the specified address
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("LISTEN error:", err)
	}
	defer conn.Close()

	fmt.Println("listening to: ", ip, ":", port)

	// listen to our connection, relaying all recieved packets to a handlePacket() go routine for proper handling
	buf := make([]byte, 1024)
	for {
		rlen, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Read error:", err)
		}

		var msg Message
		json.Unmarshal(buf[:rlen], &msg)

		// if we are to terminate our node we want to stop listening
		if msg.RPCtype == "TERMINATE_NODE" {
			fmt.Println("ohno ive been murdered")
			return
		}

		go network.handlePacket(msg)
	}
}

func (network *Network) updateBucket(sender Contact) {

	bucket := *network.Kademlia.RoutingTable.buckets[network.Kademlia.RoutingTable.getBucketIndex(network.Kademlia.RoutingTable.me.ID)]

	if bucket.Len() <= bucketSize {
		// if the bucket in nonfull we just add the new contact
		fmt.Println("contact updated in non-full bucket")
		network.Kademlia.RoutingTable.AddContact(sender)
	} else {
		// bucket is full but sender might still be in the bucket
		closestContact := network.Kademlia.LookupContact(&sender)

		// find closest contact has 0 distance means we are already in the bucket
		if closestContact[0].distance.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
			fmt.Println("existing contact updated in full bucket")
			network.Kademlia.RoutingTable.AddContact(sender) // should move us to tail of bucket
		} else {
			// ping buckets head to see if alive
			response, _ := network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			if response != nil {
				// if alive we drop the new contact
				fmt.Println("new contact dropped since bucket is alive")
				return
			} else {
				// if no response we remove dead contact and replace it with the new sender
				network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
				network.Kademlia.RoutingTable.AddContact(sender)
				fmt.Println("new contact replaced dead node in full bucket")
			}
		}
	}
}

/*
Handles the incoming packet, will do different things according to value of msg.RPCtype.
*/
func (network *Network) handlePacket(msg Message) {

	switch msg.RPCtype {
	case "PING":
		fmt.Println("you can ping, you can jive, having the time of your life")

		// add sender to my bucket
		network.updateBucket(msg.Sender)

		// send ack back
		ack := Message{
			Message: []byte("ping PONG, i hear you!"),
			RPCtype: "PING_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "PING_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println(msg.Message)

	case "FIND_CONTACT":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("find me, find me, find me a contact after midnight")
		/*
			TODO:
			do more kademlia stuff, what does a node do when it recieves a find_contact message?
		*/

		// send ack back
		ack := Message{
			Message: []byte("find contact acknowledged"),
			RPCtype: "FIND_CONTACT_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			// PROLLY MORE STUFF
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_CONTACT_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		/*
			TODO:
		*/

	case "FIND_DATA":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("data, data, data, must be funny in the rich mans world")
		/*
			TODO:
		*/

		ack := Message{
			Message: []byte("find data acknowledged"),
			RPCtype: "FIND_DATA_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			// PROLLY MORE STUFF
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_DATA_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		/*
			TODO:
			print data
			print nodeID of node contain data
		*/
		fmt.Println("I found the data you were looking for:", msg.Data)

	case "STORE":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("the winner stores it all, the loser has to fall")

		hash := network.Kademlia.Store(msg.Data)

		ack := Message{
			Message: []byte("store acknowledged"),
			RPCtype: "STORE_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			Hash:    hash,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "STORE_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("the stored data has been stored with the hash: ", msg.Hash)

		fmt.Println(msg.Hash)
	default:
		fmt.Println("oh no unknown message type recieved")
	}

}

// sends a message and returns its response if any... i hope
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
		fmt.Println("i heard something...")
		return buff[:len], nil
	} else {
		return nil, err
	}
}

// this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) ([]byte, error) {
	msg := Message{
		Message: []byte("PING pong! this is a PING message!"),
		RPCtype: "PING",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	return network.sendMessage(contact.Address, msg)
}

/*
Will tell the Listener to terminate itself.
*/
func (network *Network) SendTerminateNodeMessage() {
	msg := Message{
		Message: []byte("At Terminate Napoleon did surrender"),
		RPCtype: "TERMINATE_NODE",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	network.sendMessage(network.Kademlia.RoutingTable.me.Address, msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	msg := Message{
		Message:      []byte("greetings traveler! this is a FIND_CONTACT message!"),
		RPCtype:      "FIND_CONTACT",
		Sender:       network.Kademlia.RoutingTable.me,
		QueryContact: contact,
	}

	alphaClosest := network.Kademlia.AlphaClosest(contact.ID, network.alpha)
	// closest := alphaClosest[0] // somewhere we want to store the contact closest to queryContact we have seen yet, question is where

	for i := 0; i <= network.alpha && i < len(alphaClosest); i++ {
		fmt.Println("i:", i)
		network.sendMessage(alphaClosest[i].Address, msg)
	}
}

/*
A FIND_VALUE RPC includes a B=160-bit key. If a corresponding value is present on the recipient, the associated data is returned.
Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned.

This is a primitive operation, not an iterative one.
*/
func (network *Network) SendFindDataMessage(hash string) { // Emma needs this to print the data and the node containing the data
	msg := Message{
		Message: []byte("greetings traveler! this is a FIND_CONTACT message!"),
		RPCtype: "FIND_CONTACT",
		Sender:  network.Kademlia.RoutingTable.me,
		Hash:    hash,
	}

	// how do i find which node to send the message to?
	// kademlia stuff i guess

	network.sendMessage("contact.Address", msg)
}

/*
The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data
and make it available for later retrieval by that key.

This is a primitive operation, not an iterative one.
*/
func (network *Network) SendStoreMessage(data []byte) { // prints hash when handling response

	// find which node we want to store the data in
	// we do this by hashing the data and finding the node closest to the value of the hash?
	// hashID := network.Kademlia.getHashID(data)
	// closestNodes := network.FindClosesetNodes(hashID) // list

	// and then tell closest node to actually store it
	// msg := Message{
	// 	Message: []byte("this is a STORE message!"),
	// 	RPCtype: "STORE",
	// 	Sender:  network.Kademlia.RoutingTable.me,
	// 	Data:    data,
	// 	Hash:    hashID,
	// }
	// network.sendMessage(closestNodes[0].Address, msg)
}
