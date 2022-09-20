package d7024e

import (
	"testing"
)

func TestKademlia(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	me := NewContact(id, "327")
	//fmt.Println("meContact: ", me.String())
	table := NewRoutingTable(me)
	//kademlia := lib.Kademlia{routingTable: table}
	kademlia := Kademlia{RoutingTable: table, K: 3, Data: make(map[string][]byte)}

	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "666")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "420")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "069")
	kademlia.RoutingTable.AddContact(c1)
	kademlia.RoutingTable.AddContact(c2)
	kademlia.RoutingTable.AddContact(c3)

	//c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "101")

	result := kademlia.LookupContact(&c1)

	if result[0].String() != c1.String() {
		t.Error("got ", result[0].String(), "want ", c1.String())
	}

	kademlia.RemoveContact(&c1)

	result = kademlia.LookupContact(&c1)

	if result[0].String() == c1.String() {
		t.Error("got ", result[0].String(), "want ", c1.String())
	}

	data1 := []byte{123, 160, 161, 255, 79, 101}
	data2 := []byte{99, 01}

	hash1 := kademlia.Store(data1)
	hash2 := kademlia.Store(data2)

	_ = hash1
	_ = hash2

	dataResult, _, _ := kademlia.LookupData(hash1)
	if string(dataResult) != string(data1) {
		t.Error("got ", string(dataResult), "want ", string(data1))
	}

	dataResult, _, _ = kademlia.LookupData(hash2)
	if string(dataResult) != string(data2) {
		t.Error("got ", string(dataResult), "want ", string(data2))
	}

}
