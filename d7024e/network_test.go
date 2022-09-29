package d7024e

import (
	"fmt"
	"testing"
)

func TestNetwork(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8003"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8004"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8005"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8006"))
	k := 4
	alpha := 3
	kademlia := Kademlia{rt, k, make(map[string][]byte)}
	network := Network{kademlia, alpha, make(chan Message, alpha)}

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

}
