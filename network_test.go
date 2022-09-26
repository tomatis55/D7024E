package main

func TestNetwork() {
	alpha := 3
	k := 4
	me := NewContact(NewKademliaID("1000000000000000000000000000000000000001"), "127.0.0.1:80")
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), k, make(map[string][]byte)}, alpha, make(chan Message, alpha)}

	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:81")
	NodeNetwork.updateBucket(contact)
}
