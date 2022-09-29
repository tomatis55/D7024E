package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	K            int
	Data         map[string][]byte
}

// Finds the closest contacts that this node know of in respect to the target contact
func (kademlia *Kademlia) LookupContact(target *Contact) ContactCandidates {
	contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.K)}
	return contacts
}

// Checks if the node can find the data, if so it will return it, otherwise it will return
// the closest contacts to the data hash
func (kademlia *Kademlia) LookupData(encodedHash string) ([]byte, ContactCandidates, bool) {
	data, ok := kademlia.Data[encodedHash]
	if !ok {
		hashID := NewKademliaID(encodedHash)
		contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(hashID, kademlia.K)}
		return nil, contacts, ok
	}
	return data, ContactCandidates{}, ok
}

// Generates a hash of the data, encodes it to a hexadecimal string and returns it
func (kademlia *Kademlia) GetHash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	generatedHash := hasher.Sum(nil)
	encodedHash := hex.EncodeToString(generatedHash)
	return encodedHash
}

// Returns a KademliaID object based of an encoded hash,
// this is needed for distance calculation of between data hashes
func (kademlia *Kademlia) GetHashID(encodedHash string) KademliaID {
	return *NewKademliaID(encodedHash)
}

// Stores the data byte array by generating a hash for it and storing it in a map
// with the hash being the key and data being the value.
// Returns the encoded hash.
func (kademlia *Kademlia) Store(data []byte) string {
	encodedHash := kademlia.GetHash(data)
	kademlia.Data[encodedHash] = data
	return encodedHash
}

// Find and removes a specific contact from this node buckets
func (kademlia *Kademlia) RemoveContact(contact *Contact) {
	bucketIndex := kademlia.RoutingTable.getBucketIndex(contact.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex]
	bucket.Remove(contact)
}

func (kademlia *Kademlia) AlphaClosest(id *KademliaID, alpha int) ContactCandidates {
	return ContactCandidates{kademlia.RoutingTable.FindClosestContacts(id, alpha)}
}

func (kademlia *Kademlia) GetAllContacts() ContactCandidates {
	contacts := ContactCandidates{}
	for _, bucket := range kademlia.RoutingTable.buckets {
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			contacts.AddOne(e.Value.(Contact))
		}
	}
	return contacts
}
