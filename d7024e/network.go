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
	Get          bool // Forgive me god
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

		// if we are to terminate our node we want to stop listening
		if msg.RPCtype == "TERMINATE_NODE" {
			return
		}

		go network.handlePacket(msg)
	}
}

/*
 */

func (network *Network) updateBucket(sender Contact) {
	sender.CalcDistance(network.Kademlia.RoutingTable.me.ID)

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
			go func() {
				network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			}()

			select {
			case _ = <-network.Channel:
				// 	// if alive we drop the new contact
				return

			case <-time.After(3 * time.Second):
				// 	// if no response we remove dead contact and replace it with the new sender
				network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
				network.Kademlia.RoutingTable.AddContact(sender)
			}
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
Handles the incoming packet, will do different things according to value of msg.RPCtype.
*/
func (network *Network) handlePacket(msg Message) {

	switch msg.RPCtype {
	case "PING":
		fmt.Println("GOT PING")

		// add sender to my bucket
		network.updateBucket(msg.Sender)

		// send ack back
		ack := Message{
			RPCtype: "PING_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "PING_ACK":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println(string("GOT PONG"))

	case "FIND_CONTACT":
		/*
			TODO:
			do more kademlia stuff, what does a node do when it recieves a find_contact message?
		*/
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		contacts := network.Kademlia.LookupContact(msg.QueryContact)

		fmt.Println("Contacts from FIND_CONTACT:", contacts)

		// send ack back
		ack := Message{
			RPCtype:  "FIND_CONTACT_ACK",
			Sender:   network.Kademlia.RoutingTable.me,
			Contacts: contacts,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_CONTACT_ACK":
		fmt.Println("GOT ", msg.RPCtype, " MESSAGE")
		/*
			TODO:
		*/
		// add sender to my bucket
		//fmt.Println("In find_contact_ack")
		network.updateBucket(msg.Sender)
		network.Channel <- msg
		//fmt.Println("contact:", msg.Sender.ID, "found, buckets updated")

	case "FIND_DATA":
		/*
			if we recieve this message its because someone found that we probably contain the data someone is asking for
			now we just want to send the data back, this should be doable by using the hash as a key in the datamap in kademlia
		*/
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("GOT ", msg.RPCtype, " MESSAGE")
		// contacts := network.Kademlia.LookupContact(msg.QueryContact)

		// recover and return the data
		data, contacts, _ := network.Kademlia.LookupData(msg.Hash)

		ack := Message{
			RPCtype:  "FIND_DATA_ACK",
			Sender:   network.Kademlia.RoutingTable.me,
			Data:     data,
			Contacts: contacts,
			Get:      msg.Get,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_DATA_ACK":
		// add sender to my bucket
		fmt.Println("GOT ", msg.RPCtype, " MESSAGE")
		network.updateBucket(msg.Sender)
		network.Channel <- msg

		if msg.Get { // I hate this solution, someone please change this
			if msg.Data != nil {
				fmt.Println("I found the data you were looking for:", string(msg.Data))
				fmt.Println("in the node:                          ", msg.Sender.ID)
			} else {
				fmt.Println("no data exist at provided hash")
			}
		}

	case "STORE":
		// add sender to my bucket
		network.updateBucket(msg.Sender)
		fmt.Println("GOT ", msg.RPCtype, " MESSAGE")

		hash := network.Kademlia.Store(msg.Data)

		ack := Message{
			RPCtype: "STORE_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			Hash:    hash,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "STORE_ACK":
		// add sender to my bucket
		fmt.Println("GOT ", msg.RPCtype, " MESSAGE")
		network.updateBucket(msg.Sender)
		network.Channel <- msg
		fmt.Println("the stored data has been stored with the hash: ", msg.Hash)

	default:
		fmt.Println("Unknown message type recieved")
	}

}

// sends a message and returns its response if any... i hope
func (network *Network) sendMessage(addr string, msg Message) {
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

	fmt.Println("SENDING ", msg.RPCtype, " MESSAGE")

	if err != nil {
		fmt.Println("Write error:", err)
	}
}

// this function will send a ping message to a contact!
func (network *Network) SendPingMessage(contact *Contact) {
	msg := Message{
		RPCtype: "PING",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	network.sendMessage(contact.Address, msg)
}

/*
Will tell the Listener to terminate itself.
*/
func (network *Network) SendTerminateNodeMessage() {
	msg := Message{
		RPCtype: "TERMINATE_NODE",
		Sender:  network.Kademlia.RoutingTable.me,
	}
	network.sendMessage(network.Kademlia.RoutingTable.me.Address, msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {

	fmt.Println("Sending a find contact message to: ", contact.Address)

	msg := Message{
		RPCtype:      "FIND_CONTACT", // basically ive found you as a contact pls add me to your bucket.
		Sender:       network.Kademlia.RoutingTable.me,
		QueryContact: contact,
	}

	network.FindClosestNodes(msg)

}

/*
A FIND_VALUE RPC includes a B=160-bit key. If a corresponding value is present on the recipient, the associated data is returned.
Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned.
This is a primitive operation, not an iterative one.
*/
func (network *Network) SendFindDataMessage(hash string) { // Emma needs this to print the data and the node containing the data

	data, _, ok := network.Kademlia.LookupData(hash)
	if ok {
		fmt.Println("I found the data you were looking for:", string(data))
		fmt.Println("in the node:                          ", network.Kademlia.RoutingTable.me.ID)
	} else {
		msg := Message{
			RPCtype: "FIND_DATA",
			Sender:  network.Kademlia.RoutingTable.me,
			Hash:    hash,
			Get:     true,
		}
		network.FindClosestNodes(msg)
	}

}

/*
The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data
and make it available for later retrieval by that key.
This is a primitive operation, not an iterative one.
*/
func (network *Network) SendStoreMessage(data []byte) { // prints hash when handling response
	// msg := Message{
	// 	RPCtype: "STORE",
	// 	Sender:  network.Kademlia.RoutingTable.me,
	// 	Data:    data,
	// }

	// find which node we want to store the data in
	// we do this by hashing the data and finding the node closest to the value of the hash?
	hash := network.Kademlia.GetHash(data)

	findDataMessage := Message{
		RPCtype: "FIND_DATA",
		Sender:  network.Kademlia.RoutingTable.me,
		Data:    data,
		Hash:    hash,
		Get:     false,
	}
	contacts := network.FindClosestNodes(findDataMessage) // list

	storeMessage := Message{
		RPCtype: "STORE",
		Sender:  network.Kademlia.RoutingTable.me,
		Data:    data,
		Hash:    hash,
	}

	network.sendMessage(contacts.contacts[0].Address, storeMessage)

	// and then tell closest node to actually store it
}

func (network *Network) FindClosestNodes(msg Message) ContactCandidates {

	// Clear the channel buffer
	for len(network.Channel) > 0 {
		<-network.Channel
	}

	var id KademliaID
	switch msg.RPCtype {
	case "FIND_DATA":
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

	//  From the shortlist it selects another alpha contacts.
	//  The only condition for this selection is that they have not already been contacted.
	//  Once again a FIND_* RPC is sent to each in parallel.
	// Each such parallel search updates closestNode, the closest node seen so far.
	// The sequence of parallel searches is continued until either no node in the sets returned
	// is closer than the closest node already seen or the initiating node has accumulated k probed
	// and known to be active contacts.
	// If a cycle doesn't find a closer node, if closestNode is unchanged,
	// then the initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.
	for {
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
					break
				}

			}

		}

	}

	shortList.Sort()

	return shortList

}
