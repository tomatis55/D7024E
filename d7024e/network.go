package d7024e

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
	Alpha    int
	Channel  chan Message // make(chan Message, alpha)
}

type Message struct {
	RPCtype      string // PING, PING_ACK, FIND_CONTACT, FIND_CONTACT_ACK, FIND_DATA, FIND_DATA_ACK, STORE, STORE_ACK
	Sender       Contact
	Message      []byte
	QueryContact *Contact
	Hash         string
	Data         []byte
	Contacts     []Contact

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

		go network.handlePacket(msg)
	}
}

/*
 */

func (network *Network) updateBucket(sender Contact) {

	bucket := *network.Kademlia.RoutingTable.buckets[network.Kademlia.RoutingTable.getBucketIndex(network.Kademlia.RoutingTable.me.ID)]

	if bucket.Len() <= bucketSize {
		// if the bucket in nonfull we just add the new contact
		network.Kademlia.RoutingTable.AddContact(sender)
	} else {
		// bucket is full but sender might still be in the bucket
		closestContact := network.Kademlia.LookupContact(&sender)

		// find closest contact has 0 distance means we are already in the bucket
		if closestContact[0].distance.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
			network.Kademlia.RoutingTable.AddContact(sender) // should move us to tail of bucket
		} else {
			// ping buckets head to see if alive
			response, _ := network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			if response != nil {
				// if alive we drop the new contact
				return
			} else {
				// if no response we remove dead contact and replace it with the new sender
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

		network.Channel <- msg

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
		*/

	case "STORE":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("the winner stores it all, the loser has to fall")
		/*
			TODO:
		*/

		ack := Message{
			Message: []byte("store acknowledged"),
			RPCtype: "STORE_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			// PROLLY MORE STUFF
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "STORE_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		/*
			TODO:
		*/
	default:
		fmt.Println("oh no unknown message type recieved")
	}

}

/*
 */
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

/*
 */
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
 */

func (network *Network) SendFindContactMessage(contact *Contact) {
	msg := Message{
		Message:      []byte("greetings traveler! this is a FIND_CONTACT message!"),
		RPCtype:      "FIND_CONTACT",
		Sender:       network.Kademlia.RoutingTable.me,
		QueryContact: contact,
	}

	alphaClosest := network.Kademlia.AlphaClosest(contact.ID, network.Alpha)
	// closest := alphaClosest[0] // somewhere we want to store the contact closest to queryContact we have seen yet, question is where

	fmt.Println("len(alphaClosest):", len(alphaClosest))

	for i := 0; i <= network.Alpha && i < len(alphaClosest); i++ {
		network.sendMessage(alphaClosest[i].Address, msg)
	}
}

/*
 */

func (network *Network) SendFindDataMessage(hash string) ([]byte, Contact, error) { // Emma needs this to return the data and the node containing the data
	/*
		A FIND_VALUE RPC includes a B=160-bit key. If a corresponding value is present on the recipient, the associated data is returned.
		Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned.
		This is a primitive operation, not an iterative one.
	*/
	msg := Message{
		Message: []byte("greetings traveler! this is a FIND_VALUE message!"),
		RPCtype: "FIND_VALUE",
		Sender:  network.Kademlia.RoutingTable.me,
		Hash:    hash,
	}

	// how do i find which node to send the message to?
	// kademlia stuff i guess

	//hashID := network.Kademlia.GetHashID(hash)

	network.FindClosestNodes(msg)

	network.sendMessage("contact.Address", msg)
	return nil, Contact{}, nil
}

/*
 */

func (network *Network) SendStoreMessage(data []byte) (string, error) { // returns hash
	/*
		The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data
		and make it available for later retrieval by that key.
		This is a primitive operation, not an iterative one.
	*/
	msg := Message{
		Message: []byte("this is a STORE message!"),
		RPCtype: "STORE",
		Sender:  network.Kademlia.RoutingTable.me,
		Data:    data,
	}

	// how do i find which node to send the message to?
	// kademlia stuff and hashing i guess, mleh

	network.sendMessage("contact.Address", msg)

	return "", nil
}

func (network *Network) FindClosestNodes(msg Message) {
	var id KademliaID
	switch msg.RPCtype {
	case "FIND_VALUE":
		id = network.Kademlia.GetHashID(msg.Hash)
	case "FIND_CONTACT":
		id = *msg.Sender.ID
	}

	//	The search begins by selecting alpha contacts from the non-empty k-bucket closest to the bucket appropriate to the key being searched on.
	//	If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	//	The contact closest to the target key, closestNode, is noted.
	alphaClosest := network.Kademlia.AlphaClosest(&id, network.Alpha)
	closestNode := alphaClosest[0]

	_ = closestNode

	//	The first alpha contacts selected are used to create a shortlist for the search.
	shortList := make([]Contact, network.Kademlia.K)
	for i := 0; i < len(alphaClosest) && i < len(shortList); i++ {
		shortList[i] = alphaClosest[i]
	}

	//	The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the shortlist.
	//	Each contact, if it is live, should normally return k triples.
	//	If any of the alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.
	for _, x := range shortList {
		go func() {
			network.sendMessage(x.Address, msg)

		}()

	}

	for range shortList {
		closestContacts := <-network.Channel
		_ = closestContacts
	}

	//TODO
}
