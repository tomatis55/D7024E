package main

import (
	. "d7024e"
	"fmt"
)

// Testing kademlia
func main() {
	// "Tetrahydrocannabinol" in hex, needs to be exactly 20 characters
	// Use fmt.Printf("%s\n", decoded)
	//id := NewKademliaID("5465747261687964726F63616E6E6162696E6F6C")
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

	fmt.Println("Contact searching for: ", c1.String())
	//fmt.Println(result)
	for _, c := range result {
		fmt.Println(c.String())
	}

	fmt.Println()

	data1 := []byte{123, 160, 161, 255, 79, 101}
	data2 := []byte{99, 01}

	fmt.Println(data1)
	fmt.Println(data2)

	hash1 := kademlia.Store(data1)
	hash2 := kademlia.Store(data2)

	_ = hash1
	_ = hash2
	fmt.Println(hash1)
	fmt.Println(hash2)

	dataResult, _, ok := kademlia.LookupData(hash1)
	if ok {
		fmt.Println("Looking for data: ", data1, ", using hash: ", hash1, " Got: ", dataResult)
	}

	dataResult, _, ok = kademlia.LookupData(hash2)
	if ok {
		fmt.Println("Looking for data: ", data2, ", using hash: ", hash2, " Got: ", dataResult)
	}

	fmt.Println()

	//fmt.Printf("%s\n", *kademlia.routingTable.me.ID) // Extracts the original id string
}
