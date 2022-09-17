package d7024e

import (
	"crypto/sha1"
	"encoding/base64"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	K            int
	Data         map[string][]byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.K)
	return contacts
}

func (kademlia *Kademlia) LookupData(hash string) ([]byte, bool) {
	data, ok := kademlia.Data[hash]
	return data, ok
}

func (kademlia *Kademlia) Store(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	generatedHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	kademlia.Data[generatedHash] = data
	return generatedHash
}

func (kademlia *Kademlia) RemoveContact(contact *Contact) {
	bucketIndex := kademlia.RoutingTable.getBucketIndex(contact.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex]

	bucket.Remove(contact)
}

func (kademlia *Kademlia) AlphaClosest(contact *Contact, alpha int) []Contact {
	kClosestContacts := kademlia.LookupContact(contact)
	count := 0
	if len(kClosestContacts) < alpha {
		alphaClosest := make([]Contact, alpha)
		for _, x := range kademlia.RoutingTable.buckets {
			for e := x.list.Front(); e != nil; e = e.Next() {
				alphaClosest = append(alphaClosest, e.Value.(Contact))
				count++
				if count == alpha {
					return alphaClosest
				}
			}

		}

	} else {
		return kClosestContacts[0:alpha]
	}

	return kClosestContacts
}
