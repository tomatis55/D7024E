package d7024e

import (
	"testing"
)

func TestContactCandidates(t *testing.T) {
	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "01")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "02")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "03")
	c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "04")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "05")
	c6 := NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "06")

	temp := []Contact{c1, c2, c3, c4, c5, c6}
	contactTestList := ContactCandidates{temp}

	if !contactTestList.Contains(c1) {
		t.Error("got ", !contactTestList.Contains(c1), "want ", contactTestList.Contains(c1))
	}

	allContacts := contactTestList.GetContacts(100)

	for _, x := range allContacts {
		if !contactTestList.Contains(x) {
			t.Error("doesn't contain all contacts")
		}
	}

	contactTestList.Remove(c1)

	if contactTestList.Contains(c1) {
		t.Error("got ", contactTestList.Contains(c1), "want ", !contactTestList.Contains(c1))
	}

	if contactTestList.Index(0).ID.Equals(c1.ID) {
		t.Error("got ", c1.String(), "want ", c2.String())
	}

}
