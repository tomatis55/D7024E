package main

import "fmt"

// Testing kademlia
func main() {

	// "Tetrahydrocannabinol" in hex, needs to be exactly 20 characters
	// Use fmt.Printf("%s\n", decoded)
	id := NewKademliaID("5465747261687964726F63616E6E6162696E6F6C")
	d := id.CalcDistance(id)
	fmt.Println(d)
	me := Contact{ID: id, Address: "123", distance: d}
	table := NewRoutingTable(me)
	kademlia := Kademlia{routingTable: table}
	fmt.Println(*kademlia.routingTable.buckets[0])
	fmt.Println("AAWDAW")

	//fmt.Printf("%s\n", *kademlia.routingTable.me.ID) // Extracts the original id string
	//routingTable := RoutingTable()
	//k := kademlia(){}
}
