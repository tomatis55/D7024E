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
	Contacts     ContactCandidates

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
		if closestContact.contacts[0].distance.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
			network.Kademlia.RoutingTable.AddContact(sender) // should move us to tail of bucket
		} else {
			// TODO TODO TODO

			// // ping buckets head to see if alive
			// response, _ := network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			// if response != nil {
			// 	// if alive we drop the new contact
			// 	return
			// } else {
			// 	// if no response we remove dead contact and replace it with the new sender
			// 	network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
			// 	network.Kademlia.RoutingTable.AddContact(sender)
			// }
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
		fmt.Println(string(msg.Message))

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

		network.Channel <- msg
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
func (network *Network) sendMessage(addr string, msg Message) error {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("DIAL error:", err)
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Second))
	defer conn.Close()

	marshalled_msg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return err
	}
	_, err = conn.Write(marshalled_msg)
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}
	return err

	// buff := make([]byte, 1024)
	// len, err := conn.Read(buff)
	// //fmt.Println("conn.Read(buff): ", "len: ", len, "	", "err: ", err, "	", "buff: ", buff)
	// if err == nil {
	// 	fmt.Println("i heard something...")
	// 	return buff[:len], nil
	// } else {
	// 	fmt.Println("I didn't hear anything...")
	// 	return nil, err
	// }
}

/*
 */
// this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) {
	msg := Message{
		Message: []byte("PING pong! this is a PING message!"),
		RPCtype: "PING",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	network.sendMessage(contact.Address, msg)
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

	fmt.Println("alphaClosest.Len():", alphaClosest.Len())

	for i := 0; i <= network.Alpha && i < alphaClosest.Len(); i++ {
		network.sendMessage(alphaClosest.contacts[i].Address, msg)
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

func (network *Network) FindClosestNodes(msg Message) ContactCandidates {
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
	closestNode := alphaClosest.contacts[0]
	nodesContacted := ContactCandidates{make([]Contact, 0)}

	_ = closestNode

	//	The first alpha contacts selected are used to create a shortlist for the search.
	shortList := ContactCandidates{alphaClosest.contacts}
	shortList.Append(alphaClosest.contacts)
	nodesContacted.Append(alphaClosest.contacts)

	//	The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the shortlist.
	//	Each contact, if it is live, should normally return k triples.
	//	If any of the alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.

	messageList := make([]Message, network.Alpha)
	for _, x := range alphaClosest.contacts {
		go func() {
			network.sendMessage(x.Address, msg)
		}()
		nodesContacted.AddOne(x)

		select {
		case res := <-network.Channel:
			messageList = append(messageList, res)

		case <-time.After(3 * time.Second):
			shortList.Remove(x)
		}

	}

	//  The node then fills the shortlist with contacts from the replies received.
	//  These are those closest to the target.
	for _, message := range messageList {
		for _, x := range message.Contacts.contacts {
			if !nodesContacted.Contains(x) {
				shortList.AddOne(x)
			}

		}
	}
	shortList.Sort()
	closestNode = shortList.contacts[0]

	//  From the shortlist it selects another alpha contacts.
	//  The only condition for this selection is that they have not already been contacted.
	//  Once again a FIND_* RPC is sent to each in parallel.
	for {
		for i := 0; i < network.Alpha; i++ {
			for _, x := range shortList.contacts {
				if !nodesContacted.Contains(x) {
					go func() {
						network.sendMessage(x.Address, msg)
					}()

					nodesContacted.AddOne(x)
					select {
					case res := <-network.Channel:
						messageList = append(messageList, res)

					case <-time.After(3 * time.Second):
						shortList.Remove(x)
					}
				}
				break

			}

		}

		// Each such parallel search updates closestNode, the closest node seen so far.
		// The sequence of parallel searches is continued until either no node in the sets returned
		// is closer than the closest node already seen or the initiating node has accumulated k probed
		// and known to be active contacts.
		//
		// If a cycle doesn't find a closer node, if closestNode is unchanged,
		// then the initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.
		shortList.Sort()
		if closestNode == shortList.contacts[0] {
			break
		} else {
			closestNode = shortList.contacts[0]
		}

		kProbed := false
		count := 0
		if shortList.Len() >= network.Kademlia.K {
			for _, x := range shortList.contacts {
				if nodesContacted.Contains(x) {
					count++
				}
				if count == network.Kademlia.K {
					kProbed = true
				}
			}
		}
		if kProbed {
			break
		}

	}

	return shortList

}
