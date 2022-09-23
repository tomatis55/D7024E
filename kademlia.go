package main

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	K            int
	Data         map[string][]byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) ContactCandidates {
	contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.K)}
	return contacts
}

func (kademlia *Kademlia) LookupData(encodedHash string) ([]byte, ContactCandidates, bool) {

	data, ok := kademlia.Data[encodedHash]
	if !ok {
		hashID := NewKademliaID(encodedHash)
		contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(hashID, kademlia.K)}
		return nil, contacts, ok

	}
	return data, ContactCandidates{}, ok
}

func (kademlia *Kademlia) GetHash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	generatedHash := hasher.Sum(nil)
	encodedHash := hex.EncodeToString(generatedHash)
	return encodedHash
}

func (kademlia *Kademlia) GetHashID(encodedHash string) KademliaID {
	return *NewKademliaID(encodedHash)
}

func (kademlia *Kademlia) Store(data []byte) string {
	encodedHash := kademlia.GetHash(data)
	kademlia.Data[encodedHash] = data
	return encodedHash
}

func (kademlia *Kademlia) RemoveContact(contact *Contact) {
	bucketIndex := kademlia.RoutingTable.getBucketIndex(contact.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex]

	bucket.Remove(contact)
}

func (kademlia *Kademlia) AlphaClosest(id *KademliaID, alpha int) ContactCandidates {
	kClosestContacts := kademlia.RoutingTable.FindClosestContacts(id, kademlia.K)
	count := 0
	if len(kClosestContacts) < alpha {
		alphaClosest := make([]Contact, alpha)
		for _, x := range kademlia.RoutingTable.buckets {
			for e := x.list.Front(); e != nil; e = e.Next() {
				alphaClosest = append(alphaClosest, e.Value.(Contact))
				count++
				if count == alpha {
					contacts := ContactCandidates{alphaClosest}
					contacts.Sort()
					return contacts
				}
			}

		}

	} else {
		contacts := ContactCandidates{kClosestContacts[0:alpha]}
		contacts.Sort()
		return contacts
	}
	contacts := ContactCandidates{kClosestContacts}
	contacts.Sort()
	return contacts
}
