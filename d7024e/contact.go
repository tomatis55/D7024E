package d7024e

import (
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	if contact.distance != nil {
		return fmt.Sprintf(`contact("%s", %s, %s")`, contact.ID, contact.Address, contact.distance.String())
	} else {
		return fmt.Sprintf(`contact("%s", %s, <nil>")`, contact.ID, contact.Address)
	}

}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	contacts []Contact
}

// Adds one contact to the contacts list
func (candidates *ContactCandidates) AddOne(contact Contact) {
	candidates.contacts = append(candidates.contacts, contact)
}

// Checks if the contactCandidates contains a contact
func (candidates *ContactCandidates) Contains(contact Contact) bool {
	contains := false
	for _, x := range candidates.contacts {
		if x.ID.Equals(contact.ID) {
			contains = true
		}
	}
	return contains
}

// Removes a contact from the internal contacts list
func (candidates *ContactCandidates) Remove(contact Contact) {
	for i, x := range candidates.contacts {
		if x.ID.Equals(contact.ID) {
			copy(candidates.contacts[i:], candidates.contacts[i+1:])               // Shift a[i+1:] left one index.
			candidates.contacts[len(candidates.contacts)-1] = Contact{}            // Erase last element (write zero value).
			candidates.contacts = candidates.contacts[:len(candidates.contacts)-1] // Truncate slice.
			return
		}
	}
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	if count > candidates.Len() {
		return candidates.contacts
	} else {
		return candidates.contacts[:count]
	}

}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}
