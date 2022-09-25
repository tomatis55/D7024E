package main

import (
	"time"
)

func TestNetwork() {
	alpha := 3
	k := 4
	me := NewContact(NewKademliaID("1000000000000000000000000000000000000001"), "127.0.0.1:80")
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), k, make(map[string][]byte)}, alpha, make(chan Message, alpha)}

	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:81")
	NodeNetwork.updateBucket(contact)
}

func (network *Network) ddd(sender Contact) {

	sender.CalcDistance(network.Kademlia.RoutingTable.me.ID)
	bucket := *network.Kademlia.RoutingTable.buckets[network.Kademlia.RoutingTable.getBucketIndex(network.Kademlia.RoutingTable.me.ID)]

	if bucket.Len() <= bucketSize {
		// if the bucket in nonfull we just add the new contact
		network.Kademlia.RoutingTable.AddContact(sender)
	} else {
		// bucket is full but sender might still be in the bucket
		closestContact := network.Kademlia.LookupContact(&sender)

		// find closest contact has 0 distance means we are already in the bucket
		if closestContact.contacts[0].distance.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
			network.Kademlia.RoutingTable.AddContact(sender) // should move us to tail of bucket
		} else {
			// TODO TODO TODO
			// ping buckets head to see if alive
			go func() {
				network.SendPingMessage(bucket.list.Front().Value.(*Contact))
			}()

			select {
			case _ = <-network.Channel:
				// if alive we drop the new contact
				return

			case <-time.After(3 * time.Second):
				// if no response we remove dead contact and replace it with the new sender
				network.Kademlia.RemoveContact(bucket.list.Front().Value.(*Contact))
				network.Kademlia.RoutingTable.AddContact(sender)
			}
		}
	}
}
