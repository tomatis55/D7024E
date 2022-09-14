package d7024e

import (
	"fmt"
	"testing"
)

func TestKademlia(*testing.T) {
	fmt.Println("Hello")
	// Output: Hello
}

func TestLookupContact(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	me := NewContact(id, "327")
	fmt.Println("meContact: ", me.String())
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
		t.Error("got %d, want %d", result[0].String(), c1.String())
	}
	fmt.Println("Contact searching for: ", c1.String())
	//fmt.Println(result)
	for _, c := range result {
		fmt.Println(c.String())
	}
}
