package main

import (
	"testing"
)

func TestKademlia(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	me := NewContact(id, "00")
	//fmt.Println("meContact: ", me.String())
	table := NewRoutingTable(me)
	//kademlia := lib.Kademlia{routingTable: table}
	kademlia := NewKademlia(table, 3)

	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "01")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "02")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "03")
	c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "04")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "05")
	c6 := NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "06")

	temp := []Contact{c1, c2, c3, c4, c5, c6}
	contactTestList := ContactCandidates{temp}

	kademlia.RoutingTable.AddContact(c1)
	kademlia.RoutingTable.AddContact(c2)
	kademlia.RoutingTable.AddContact(c3)
	kademlia.RoutingTable.AddContact(c4)
	kademlia.RoutingTable.AddContact(c5)
	kademlia.RoutingTable.AddContact(c6)

	result := kademlia.LookupContact(&c6)

	if result.contacts[0].String() != c6.String() {
		t.Error("got ", result.contacts[0].String(), "want ", c1.String())
	}

	data1 := []byte{123, 160, 161, 255, 79, 101}
	data2 := []byte{99, 01}
	data3 := []byte{69, 69, 3}

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

	hash3 := kademlia.GetHash(data3)
	_, _, ok := kademlia.LookupData(hash3)
	if ok {
		t.Error("got ", ok, "want ", !ok)
	}

	hashID := kademlia.GetHashID(hash3)
	if hashID.String() != "c03e151c0105c203a0d794b19fb20352155f6c2b" {
		t.Error("got ", hashID.String(), "want ", "c03e151c0105c203a0d794b19fb20352155f6c2b")
	}

	contacts := kademlia.AlphaClosest(id, 3)
	for i := 0; i < contacts.Len(); i++ {
		if !contacts.contacts[i].ID.Equals(contactTestList.contacts[i].ID) {
			t.Error("got ", contacts.contacts[i].String(), "want ", contactTestList.contacts[i].String())
		}
	}

	contacts = kademlia.AlphaClosest(id, 4)
	for i := 0; i < contacts.Len(); i++ {
		if !contacts.contacts[i].ID.Equals(contactTestList.contacts[i].ID) {
			t.Error("got ", contacts.contacts[i].String(), "want ", contactTestList.contacts[i].String())
		}
	}

	allContacts := kademlia.GetAllContacts()
	for i := 0; i < allContacts.Len(); i++ {
		if !allContacts.contacts[i].ID.Equals(contactTestList.contacts[5-i].ID) {
			t.Error("got ", allContacts.contacts[i].String(), "want ", contactTestList.contacts[5-i].String())
		}
	}

	kademlia.RemoveContact(&c6)

	result = kademlia.LookupContact(&c6)

	if result.contacts[0].String() == c6.String() {
		t.Error("got ", result.contacts[0].String(), "want ", c1.String())
	}

}

func TestKademliaID(t *testing.T) {
	randomId1 := NewRandomKademliaID()
	randomId2 := NewRandomKademliaID()

	if randomId1.Equals(randomId2) {
		t.Error("Not randomly generated ID")
	}
}
