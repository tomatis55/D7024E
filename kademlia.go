package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

var TimeToLive time.Duration = time.Second * 30 // deathTimer

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

	go func(ttl time.Duration, hash string) {
		time.Sleep(ttl)
		fmt.Println("hello! data at hash", hash, "is expiring now")
		kademlia.Data[hash] = nil
	}(TimeToLive, encodedHash)

	return encodedHash
}

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
