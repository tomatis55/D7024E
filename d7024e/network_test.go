package d7024e

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNetwork(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1:8000"))
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

	network.updateBucket(c)
	contacts := network.Kademlia.LookupContact(&c)
	if contacts.contacts[0].ID != c.ID {
		t.Error("got ", contacts.contacts[0].ID.String(), "want ", c.ID.String())
	}

	msg := Message{}

	err := network.sendMessage("localhost:8001", msg)

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

	go network.Listen("127.0.0.1", 8000)
	contacts2 := make([]Contact, 0, 3)
	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000008"), "localhost:8008")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000009"), "localhost:8009")
	c3 := NewContact(NewKademliaID("000000000000000000000000000000000000000a"), "localhost:8010")
	contacts2 = append(contacts2, c1, c2, c3)
	msg2 := Message{RPCtype: "FIND_CONTACT_ACK", Sender: network.Kademlia.RoutingTable.me, Contacts: contacts2}
	err3 := network.sendMessage("127.0.0.1:8000", msg2)

	if err3 != nil {
		t.Error("got ", err3, "want ", nil)
	}
	for i, contact := range contacts2 {
		contact.CalcDistance(NewKademliaID("0000000000000000000000000000000000000000"))
		contacts2[i] = contact
	}

	_ = network.sendMessage("127.0.0.1:8000", msg2)

	// TODO: check if the contacts are updated correctly in the bucket

	msgReceived := <-network.MsgChannel

	if !reflect.DeepEqual(msgReceived.Contacts, contacts2) {
		t.Error("got ", msgReceived.Contacts, "want ", contacts2)
	}

	network.SendTerminateNodeMessage()
	c4 := NewContact(NewKademliaID("000000000000000000000000000000000000000b"), "localhost:8011")

	candidates := network.SendFindContactMessage(&c4)

	if candidates.Len() != 0 {
		t.Error("got ", candidates.contacts, "want {}")
	}

	candidates2 := network.SendFindDataMessage("0000000000000000000400000000000000000000")
	if candidates2.Len() != 0 {
		t.Error("got ", candidates2.contacts, "want {}")
	}

	rt2 := NewRoutingTable(NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012"))
	rt2.AddContact(NewContact(NewKademliaID("000000000000000000000000000000000000000d"), "127.0.0.1:8013"))
	rt2.AddContact(NewContact(NewKademliaID("111111111111111111111111111111111111111e"), "127.0.0.1:8014"))
	kademlia2 := Kademlia{rt2, k, make(map[string][]byte)}
	network2 := Network{kademlia2, alpha, make(chan Message, alpha), make(chan []byte)}

	rt3 := NewRoutingTable(NewContact(NewKademliaID("000000000000000000000000000000000000000d"), "127.0.0.1:8013"))
	rt3.AddContact(NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012"))
	kademlia3 := Kademlia{rt3, k, make(map[string][]byte)}
	network3 := Network{kademlia3, alpha, make(chan Message, alpha), make(chan []byte)}

	rt4 := NewRoutingTable(NewContact(NewKademliaID("111111111111111111111111111111111111111e"), "127.0.0.1:8014"))
	rt4.AddContact(NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012"))
	kademlia4 := Kademlia{rt4, k, make(map[string][]byte)}
	network4 := Network{kademlia4, alpha, make(chan Message, alpha), make(chan []byte)}

	fmt.Println("=========STARTING NEW NETWORKS=========")

	go network2.Listen("127.0.0.1", 8012)
	go network3.Listen("127.0.0.1", 8013)
	go network4.Listen("127.0.0.1", 8014)

	c0 := NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012")
	candidates3 := network3.SendFindContactMessage(&c0)
	for _, x := range candidates3.contacts {
		fmt.Println(x.String())
	}

	if !candidates3.contacts[0].ID.Equals(NewKademliaID("000000000000000000000000000000000000000c")) {
		t.Error("got ", candidates3.contacts[0].Address, "want: ", "127.0.0.1:8012")
	}

	rt5 := NewRoutingTable(NewContact(NewKademliaID("111111111111111111111111111111111111111f"), "127.0.0.1:8015"))
	rt5.AddContact(NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012"))
	rt2.AddContact(NewContact(NewKademliaID("111111111111111111111111111111111111111f"), "127.0.0.1:8015"))
	kademlia5 := Kademlia{rt5, k, make(map[string][]byte)}
	network5 := Network{kademlia5, alpha, make(chan Message, alpha), make(chan []byte)}
	go network5.Listen("127.0.0.1", 8015)

	e := NewContact(NewKademliaID("111111111111111111111111111111111111111f"), "127.0.0.1:8015")
	candidates4 := network3.SendFindContactMessage(&e)

	if !candidates4.contacts[0].ID.Equals(NewKademliaID("111111111111111111111111111111111111111f")) {
		t.Error("got ", candidates4.contacts[0].Address, "want: ", "127.0.0.1:8015")
	}

	network2.SendTerminateNodeMessage()
	network3.SendTerminateNodeMessage()
	network4.SendTerminateNodeMessage()
	network5.SendTerminateNodeMessage()
	//t.Error()
}
