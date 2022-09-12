package main

import "fmt"

// Testing kademlia
func main() {

	// "Tetrahydrocannabinol" in hex, needs to be exactly 20 characters
	// Use fmt.Printf("%s\n", decoded)
	id := NewKademliaID("5465747261687964726F63616E6E6162696E6F6C")
	d := id.CalcDistance(id)
	me := Contact{ID: id, Address: "123", distance: d}
	table := NewRoutingTable(me)
	kademlia := Kademlia{routingTable: table}
	//fmt.Printf("%s\n", *kademlia.routingTable.me.ID) // Extracts the original id string
	fmt.Println("routingTable.buckets[0]: ", *kademlia.routingTable.buckets[0])
	testID := NewRandomKademliaID()
	fmt.Println("testID: ", testID)
	testDistance := testID.CalcDistance(id)
	fmt.Println("Distance between me and testID: ", testDistance)
	testContact := Contact{ID: testID, Address: "327", distance: testDistance}
	kademlia.routingTable.AddContact(testContact)
	//fmt.Println("routingTable.buckets[0]: ", *kademlia.routingTable.buckets)
	//routingTable := RoutingTable()
	//k := kademlia(){}
}
