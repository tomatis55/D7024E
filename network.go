package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Network struct {
	Kademlia Kademlia
	Alpha    int
	Channel  chan Message
}

type Message struct {
	RPCtype      string // PING, PING_ACK, FIND_CONTACT, FIND_CONTACT_ACK, FIND_DATA, FIND_DATA_ACK, STORE, STORE_ACK
	Sender       Contact
	QueryContact *Contact
	Hash         string
	Data         []byte
	Contacts     []Contact
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
			// fmt.Println("ohno ive been murdered")
			return
		}

		go network.handlePacket(msg)
	}
}

func (network *Network) updateBucket(sender Contact) {

	// i dont want to be in my own bucket
	if network.Kademlia.RoutingTable.me.ID.Equals(sender.ID) {
		return
	}

	// fmt.Println("in update bucket, sender, before:", sender)

	sender.CalcDistance(network.Kademlia.RoutingTable.me.ID)

	// fmt.Println("in update bucket, sender, after:", sender)

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
			// ping buckets head to see if alive
			go func() {
				network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			}()

			select {
			case _ = <-network.Channel:
				// if alive we drop the new contact
				return

			case <-time.After(3 * time.Second):
				// if no response we remove dead contact and replace it with the new sender
				network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
				network.Kademlia.RoutingTable.AddContact(sender)
			}
		}
	}
}

/*
Handles the incoming packet, will do different things according to value of msg.RPCtype.
*/
func (network *Network) handlePacket(msg Message) {
	// add sender to my bucket
	network.updateBucket(msg.Sender)

	switch msg.RPCtype {
	case "PING":
		fmt.Println("you can ping, you can jive, having the time of your life")

		// send ack back
		ack := Message{
			RPCtype: "PING_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "PING_ACK":
		fmt.Println("ping PONG, i hear you!")

	case "FIND_CONTACT":
		fmt.Println("find me, find me, find me a contact after midnight")
		contacts := network.Kademlia.LookupContact(msg.QueryContact)

		// send ack back
		ack := Message{
			RPCtype:  "FIND_CONTACT_ACK",
			Sender:   network.Kademlia.RoutingTable.me,
			Contacts: contacts.contacts,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_CONTACT_ACK":

		for i, contact := range msg.Contacts {
			// fmt.Println("ID:", network.Kademlia.RoutingTable.me.ID)
			contact.CalcDistance(network.Kademlia.RoutingTable.me.ID)
			msg.Contacts[i] = contact
			network.updateBucket(contact)
			// fmt.Println("contact:", contact.ID, "found, buckets updated")
		}

		// fmt.Println("\nmsg before channel:", msg, "\n")
		network.Channel <- msg

	case "FIND_DATA":
		/*
			if we recieve this message its because someone found that we probably contain the data someone is asking for
			now we just want to send the data back, this should be doable by using the hash as a key in the datamap in kademlia
		*/
		// add sender to my bucket

		_, contacts, _ := network.Kademlia.LookupData(msg.Hash)

		fmt.Println("data, data, data, must be funny in the rich mans world")

		ack := Message{
			RPCtype:  "FIND_DATA_ACK",
			Sender:   network.Kademlia.RoutingTable.me,
			Hash:     msg.Hash,
			Contacts: contacts.contacts,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "FIND_DATA_ACK":

		for i, contact := range msg.Contacts {
			contact.CalcDistance(NewKademliaID(msg.Hash))
			msg.Contacts[i] = contact
		}

		network.Channel <- msg

	case "RECOVER_DATA":

		data, _, _ := network.Kademlia.LookupData(msg.Hash)

		// fmt.Println("Data from lookupData: ", data)

		ack := Message{
			RPCtype: "RECOVER_DATA_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			Data:    data,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "RECOVER_DATA_ACK":

		if msg.Data != nil {
			fmt.Println("I found the data you were looking for:", string(msg.Data))
			fmt.Println("in the node:", msg.Sender.ID)
		} else {
			fmt.Println("no data exist at provided hash :(")
		}

	case "STORE":
		fmt.Println("the winner stores it all, the loser has to fall")

		hash := network.Kademlia.Store(msg.Data)

		ack := Message{
			RPCtype: "STORE_ACK",
			Sender:  network.Kademlia.RoutingTable.me,
			Hash:    hash,
		}
		network.sendMessage(msg.Sender.Address, ack)

	case "STORE_ACK":
		network.Channel <- msg
		fmt.Println("the data has been stored with the hash: ", msg.Hash)

	default:
		fmt.Println("oh no unknown message type recieved")
	}

}

func (network *Network) sendMessage(addr string, msg Message) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("DIAL error:", err)
	}
	// conn.SetDeadline(time.Now().Add(time.Second))
	defer conn.Close()

	marshalled_msg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Marshal error:", err)
	}
	_, err = conn.Write(marshalled_msg)
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

func (network *Network) SendFindContactMessage(contact *Contact) ContactCandidates {

	msg := Message{
		RPCtype:      "FIND_CONTACT", // basically ive found you as a contact pls add me to your bucket.
		Sender:       network.Kademlia.RoutingTable.me,
		QueryContact: contact,
	}

	shortList := network.FindClosestNodes(msg)
	return shortList
}

/*
A FIND_VALUE RPC includes a B=160-bit key. If a corresponding value is present on the recipient, the associated data is returned.
Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned.

This is a primitive operation, not an iterative one.
*/
func (network *Network) SendFindDataMessage(hash string) {
	data, _, ok := network.Kademlia.LookupData(hash)
	if ok {
		// if data is in local node, print it
		fmt.Println("I found the data you were looking for:", string(data))
		fmt.Println("in the node:                          ", network.Kademlia.RoutingTable.me.ID)
	} else {
		findDataMessage := Message{
			RPCtype: "FIND_DATA",
			Sender:  network.Kademlia.RoutingTable.me,
			// Data:    data,
			Hash: hash,
		}
		contacts := network.FindClosestNodes(findDataMessage) // list

		recoverDataMessage := Message{
			RPCtype: "RECOVER_DATA",
			Sender:  network.Kademlia.RoutingTable.me,
			Hash:    hash,
		}
		fmt.Println("Sending recover data msg to: ", contacts.contacts[0].Address)
		network.sendMessage(contacts.contacts[0].Address, recoverDataMessage)
	}
}

/*
The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data
and make it available for later retrieval by that key.

This is a primitive operation, not an iterative one.
*/
func (network *Network) SendStoreMessage(data []byte) { // prints hash when handling response
	// find which node we want to store the data in
	// we do this by hashing the data and finding the node closest to the value of the hash?
	hash := network.Kademlia.GetHash(data)

	findDataMessage := Message{
		RPCtype: "FIND_DATA",
		Sender:  network.Kademlia.RoutingTable.me,
		Data:    data,
		Hash:    hash,
	}
	contacts := network.FindClosestNodes(findDataMessage) // list

	storeMessage := Message{
		RPCtype: "STORE",
		Sender:  network.Kademlia.RoutingTable.me,
		Data:    data,
		Hash:    hash,
	}

	distFromMe := network.Kademlia.RoutingTable.me.ID.CalcDistance(NewKademliaID(hash))
	// fmt.Println("distFromMe:", distFromMe)                          // has differing actual values
	// fmt.Println("closest contact:", contacts.contacts[0])           // is always me
	// fmt.Println("contact.distance:", contacts.contacts[0].distance) // is always 0000000000000000000000000000000000000000

	// fmt.Println("n.K.R.m.A:", network.Kademlia.RoutingTable.me.Address) // this
	// fmt.Println("c.c[0].A :", contacts.contacts[0].Address)             // and this are always same

	if distFromMe.Less(contacts.contacts[0].distance) { // never enters here
		network.sendMessage(network.Kademlia.RoutingTable.me.Address, storeMessage)
	} else { // always enters here
		network.sendMessage(contacts.contacts[0].Address, storeMessage)
	}
}

func (network *Network) FindClosestNodes(msg Message) ContactCandidates {
	var id KademliaID
	// var closestNode Contact
	switch msg.RPCtype {
	case "FIND_DATA":
		id = network.Kademlia.GetHashID(msg.Hash)
		// alphaClosest := network.Kademlia.AlphaClosest(&id, network.Alpha)
		// closestNode = alphaClosest.contacts[0]
	case "FIND_CONTACT":
		id = *msg.Sender.ID
		// closestNode = network.Kademlia.RoutingTable.me
	}

	//	The search begins by selecting alpha contacts from the non-empty k-bucket closest to the bucket appropriate to the key being searched on.
	//	If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	//	The contact closest to the target key, closestNode, is noted.
	alphaClosest := network.Kademlia.AlphaClosest(&id, network.Alpha)
	closestNode := alphaClosest.contacts[0]
	// closestNode := network.Kademlia.RoutingTable.me
	nodesContacted := ContactCandidates{make([]Contact, 0)}

	//	The first alpha contacts selected are used to create a shortlist for the search.
	shortList := ContactCandidates{alphaClosest.contacts}
	// shortList.Append(alphaClosest.contacts)
	nodesContacted.Append(alphaClosest.contacts)

	// fmt.Println("shortlist in start:\n", shortList)

	//	The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the shortlist.
	//	Each contact, if it is live, should normally return k triples.
	//	If any of the alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.

	messageList := make([]Message, network.Alpha)
	for _, x := range alphaClosest.contacts {
		for len(network.Channel) > 0 {
			<-network.Channel
		}
		go func() {
			// fmt.Println("1Sending message to: ", x.Address)
			network.sendMessage(x.Address, msg)
		}()
		nodesContacted.AddOne(x)

		select {
		case res := <-network.Channel:
			messageList = append(messageList, res)

		case <-time.After(10 * time.Second):
			// fmt.Println(x.Address)
			shortList.Remove(x)
		}
	}

	// fmt.Println("shortlist before response:\n", shortList)

	//  The node then fills the shortlist with contacts from the replies received.
	//  These are those closest to the target.
	for _, message := range messageList {
		for _, x := range message.Contacts {
			if !nodesContacted.Contains(x) && !shortList.Contains(x) {
				shortList.AddOne(x)
			}
		}
	}

	//  From the shortlist it selects another alpha contacts.
	//  The only condition for this selection is that they have not already been contacted.
	//  Once again a FIND_* RPC is sent to each in parallel.
	if shortList.Len() == 0 {
		return shortList
	}
	for {
		shortList.Sort()
		// fmt.Println("Shortlist: ", shortList)
		if closestNode.ID.Equals(shortList.contacts[0].ID) {
			break
		} else {
			closestNode = shortList.contacts[0]
		}

		for i := 0; i < network.Alpha; i++ {
			for _, x := range shortList.contacts {
				if !nodesContacted.Contains(x) {
					for len(network.Channel) > 0 {
						<-network.Channel
					}
					go func() {
						// fmt.Println("2Sending message to: ", x.Address)
						network.sendMessage(x.Address, msg)
					}()

					nodesContacted.AddOne(x)
					select {
					case res := <-network.Channel:
						// fmt.Println("res from channel:", res)
						messageList = append(messageList, res)

					case <-time.After(3 * time.Second):
						shortList.Remove(x)
					}
					break
				}
			}
		}
		for _, message := range messageList {
			for _, x := range message.Contacts {
				if !nodesContacted.Contains(x) && !shortList.Contains(x) {
					shortList.AddOne(x)
				}
			}
		}
		// Each such parallel search updates closestNode, the closest node seen so far.
		// The sequence of parallel searches is continued until either no node in the sets returned
		// is closer than the closest node already seen or the initiating node has accumulated k probed
		// and known to be active contacts.
		//
		// If a cycle doesn't find a closer node, if closestNode is unchanged,
		// then the initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.

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

	shortList.Sort()
	return shortList
}
