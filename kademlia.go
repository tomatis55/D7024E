package d7024e

type Kademlia struct {
	RoutingTable *RoutingTable
	Data         map[string]byte
	K            int
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.K)
	return contacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
