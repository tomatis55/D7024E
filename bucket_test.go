package main

import (
	"testing"
)

func TestBuckets(t *testing.T) {
	c1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "01")
	c2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "02")
	c3 := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "03")
	c4 := NewContact(NewKademliaID("0000000000000000000000000000000000000004"), "04")
	c5 := NewContact(NewKademliaID("0000000000000000000000000000000000000005"), "05")
	c6 := NewContact(NewKademliaID("0000000000000000000000000000000000000006"), "06")
	c7 := NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "07")

	temp := []Contact{c7, c6, c5, c4, c3, c2, c1}

	bucket := newBucket()

	bucket.AddContact(c1)
	bucket.AddContact(c2)
	bucket.AddContact(c3)
	bucket.AddContact(c4)
	bucket.AddContact(c5)
	bucket.AddContact(c6)
	bucket.AddContact(c7)

	i := 0
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		if !contact.ID.Equals(temp[i].ID) {
			t.Error("got ", contact.String(), "want ", temp[i].String())
		}
		i++
	}

	bucket.AddContact(c3)

	front := bucket.list.Front().Value.(Contact)
	if !front.ID.Equals(c3.ID) {
		t.Error("got ", front.String(), "want ", c3.String())
	}

}
