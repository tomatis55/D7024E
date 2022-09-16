package d7024e

import (
	"testing"
)

func TestRoutingTable(t *testing.T) {

	tables := []struct {
		ID string
	}{
		{"2111111400000000000000000000000000000000"},
		{"1111111400000000000000000000000000000000"},
		{"1111111100000000000000000000000000000000"},
		{"1111111200000000000000000000000000000000"},
		{"1111111300000000000000000000000000000000"},
		{"ffffffff00000000000000000000000000000000"},
	}

	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range tables {
		if tables[i].ID != contacts[i].ID.String() {
			t.Error("got ", contacts[i].ID.String(), "want ", tables[i].ID)
		}

	}
}
