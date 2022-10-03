package d7024e

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNetwork(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "172.17.0.1:8000"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8003"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8004"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8005"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8006"))
	k := 4
	alpha := 3
	kademlia := Kademlia{rt, k, make(map[string][]byte)}
	network := Network{kademlia, alpha, make(chan Message, alpha), make(chan []byte)}

	c := NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "localhost:8007")

	fmt.Println(network.Kademlia.LookupContact(&c))
	network.updateBucket(c)
	contacts := network.Kademlia.LookupContact(&c)
	if contacts.contacts[0].ID != c.ID {
		t.Error("got ", contacts.contacts[0].ID.String(), "want ", c.ID.String())
	}

	msg := Message{}

	err := network.sendMessage("localhost:8000", msg)

	if err != nil {
		t.Error("got ", err, "want ", nil)
	}

	err2 := network.SendPingMessage(&c)

	if err2 != nil {
		t.Error("got ", err2, "want ", nil)
	}

	// go listen , make sure main thread is active long enough
	// send find_contact_ack msg
	// check if the contacts are updated correctly
	// check the channel to see if the msg is intact

	go network.Listen("172.17.0.1", 8000)
	contacts2 := make([]Contact, 0, 3)
	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000008"), "localhost:8008")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000009"), "localhost:8009")
	c3 := NewContact(NewKademliaID("000000000000000000000000000000000000000a"), "localhost:8010")
	contacts2 = append(contacts2, c1, c2, c3)
	msg2 := Message{RPCtype: "FIND_CONTACT_ACK", Sender: network.Kademlia.RoutingTable.me, Contacts: contacts2}
	err3 := network.sendMessage("172.17.0.1:8000", msg2)

	if err3 != nil {
		t.Error("got ", err3, "want ", nil)
	}

	for i, contact := range contacts2 {
		contact.CalcDistance(NewKademliaID("0000000000000000000000000000000000000000"))
		contacts2[i] = contact
	}

	_ = network.sendMessage("172.17.0.1:8000", msg2)

	fmt.Println(contacts2)

	// TODO: check if the contacts are updated correctly in the bucket

	time.Sleep(2 * time.Second)

	fmt.Println("after sleep, before receive")
	msgReceived := <-network.MsgChannel

	fmt.Println("after sleep and receive")

	if !reflect.DeepEqual(msgReceived.Contacts, contacts2) {
		t.Error("got ", msgReceived.Contacts, "want ", contacts2)
	}

	fmt.Println("passed??? somehow??!?!?!?")

	network.SendTerminateNodeMessage()

	// test findCLosestNodes
	// expect shortlist to be empty since no one answers??
	c4 := NewContact(NewKademliaID("000000000000000000000000000000000000000b"), "172.17.0.5:8011")
	// msg3 := Message{
	// 	RPCtype:      "FIND_CONTACT",
	// 	Sender:       network.Kademlia.RoutingTable.me,
	// 	QueryContact: &c4,
	// }
	// candidates := network.FindClosestNodes(msg3)

	candidates := network.SendFindContactMessage(&c4)

	if candidates.Len() != 0 {
		t.Error("got ", candidates.contacts, "want {}")
	}

	candidates2 := network.SendFindDataMessage("0000000000000000000400000000000000000000")
	if candidates2.Len() != 0 {
		t.Error("got ", candidates2.contacts, "want {}")
	}
}
