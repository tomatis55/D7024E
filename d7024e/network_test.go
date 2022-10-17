package d7024e

import (
	"fmt"
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
	kademlia := NewKademlia(rt, k)
	network := NewNetwork(kademlia, alpha)

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
	s := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1:8000")
	msg2 := Message{RPCtype: "FIND_CONTACT_ACK", Sender: network.Kademlia.RoutingTable.me, Contacts: contacts2, QueryContact: &s}
	err3 := network.sendMessage("127.0.0.1:8000", msg2)

	if err3 != nil {
		t.Error("got ", err3, "want ", nil)
	}
	for i, contact := range contacts2 {
		contact.CalcDistance(NewKademliaID("0000000000000000000000000000000000000000"))
		contacts2[i] = contact
	}

	// TODO: check if the contacts are updated correctly in the bucket

	// msgReceived := <-network.MsgChannel

	// if !reflect.DeepEqual(msgReceived.Contacts, contacts2) {
	// 	t.Error("got ", msgReceived.Contacts, "want ", contacts2)
	// }

	network.SendTerminateNodeMessage()
	c4 := NewContact(NewKademliaID("000000000000000000000000000000000000000b"), "localhost:8011")

	candidates := network.SendFindContactMessage(&c4)

	if candidates.Len() != 0 {
		t.Error("got ", candidates.contacts, "want {}")
	}

	_, _, ok := network.SendFindDataMessage("0000000000000000000400000000000000000000")
	if ok {
		t.Error("got ", ok, "want ", !ok)
	}

	superContact := NewContact(NewKademliaID("000000000000000000000000000000000000000c"), "127.0.0.1:8012")
	d3 := NewContact(NewKademliaID("000000000000000000000000000000000000000d"), "127.0.0.1:8013")
	e4 := NewContact(NewKademliaID("111111111111111111111111111111111111111e"), "127.0.0.1:8014")

	rt2 := NewRoutingTable(superContact)
	rt2.AddContact(d3)
	rt2.AddContact(e4)
	kademlia2 := NewKademlia(rt2, k)
	network2 := NewNetwork(kademlia2, alpha)

	rt3 := NewRoutingTable(d3)
	rt3.AddContact(superContact)
	kademlia3 := NewKademlia(rt3, k)
	network3 := NewNetwork(kademlia3, alpha)

	rt4 := NewRoutingTable(e4)
	rt4.AddContact(superContact)
	kademlia4 := NewKademlia(rt4, k)
	network4 := NewNetwork(kademlia4, alpha)

	fmt.Println("=========STARTING NEW NETWORKS=========")

	go network2.Listen("127.0.0.1", 8012)
	go network3.Listen("127.0.0.1", 8013)
	go network4.Listen("127.0.0.1", 8014)

	candidates3 := network3.SendFindContactMessage(&superContact)
	for _, x := range candidates3.contacts {
		fmt.Println(x.String())
	}

	if !candidates3.contacts[0].ID.Equals(superContact.ID) {
		t.Error("got ", candidates3.contacts[0].Address, "want: ", "127.0.0.1:8012")
	}

	f5 := NewContact(NewKademliaID("111111111111111111111111111111111111111f"), "127.0.0.1:8015")
	rt5 := NewRoutingTable(f5)
	rt5.AddContact(superContact)
	rt2.AddContact(f5)
	kademlia5 := NewKademlia(rt5, k)
	network5 := NewNetwork(kademlia5, alpha)
	go network5.Listen("127.0.0.1", 8015)

	network3.SendTerminateNodeMessage()

	candidates4 := network4.SendFindContactMessage(&f5)

	fmt.Println("candidates4.contacts: ", candidates4.contacts)
	if !candidates4.contacts[0].ID.Equals(f5.ID) {
		t.Error("got ", candidates4.contacts[0].Address, "want: ", "127.0.0.1:8015")
	}
	data := []byte("Test data")
	dataContacts, hash := network2.SendStoreMessage(data)
	if !dataContacts.contacts[0].ID.Equals(superContact.ID) {
		t.Error("got ", dataContacts.contacts[0].Address, "want: ", superContact.Address)
	}

	a6 := NewContact(NewKademliaID("111111111111111111111111111111111111111a"), "127.0.0.1:8016")
	rt6 := NewRoutingTable(a6)
	rt6.AddContact(superContact)
	rt2.AddContact(a6)
	kademlia6 := NewKademlia(rt6, k)
	network6 := NewNetwork(kademlia6, alpha)
	go network6.Listen("127.0.0.1", 8016)

	dataContactID, dataRecieved, ok := network6.SendFindDataMessage(hash)
	fmt.Println(dataContactID)
	if !ok {
		t.Error("want ", !ok, "got ", ok)
	} else if dataRecieved != string(data) {
		t.Error("want ", data, "got ", dataRecieved)
	}

	network2.SendTerminateNodeMessage()
	network4.SendTerminateNodeMessage()
	network5.SendTerminateNodeMessage()
	network6.SendTerminateNodeMessage()
	//t.Error()
}
