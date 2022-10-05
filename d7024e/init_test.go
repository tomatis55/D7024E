package d7024e

import (
	"testing"
)

func TestInit(t *testing.T) {
	_ = InitalizeSuperNode("0000000000000000000000000000000000000111", "127.0.0.1", 8020)
	network := InitalizeNode("127.0.0.1", 8021, "0000000000000000000000000000000000000111", "127.0.0.1", "8020")

	if !network.Kademlia.GetAllContacts().contacts[0].ID.Equals(NewKademliaID("0000000000000000000000000000000000000111")) {
		t.Error("got ", network.Kademlia.GetAllContacts().contacts[0].ID.String(), "want ", "0000000000000000000000000000000000000111")
	}

	// if superNetwork.Kademlia.GetAllContacts().contacts[0].Address != "127.0.0.1:8021" {
	// 	t.Error("got ", superNetwork.Kademlia.GetAllContacts().contacts[0].Address, "want ", "127.0.0.1:8021")
	// }

}
