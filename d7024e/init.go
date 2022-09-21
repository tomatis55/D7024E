package d7024e

// "fmt"
// "time"

var NodeNetwork Network

func InitalizeSuperNode(id string, ip string) {
	alpha := 3
	me := NewContact(NewKademliaID(id), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}, alpha, make(chan Message)}
}

func InitalizeNode(ip string, idSuperNode string, ipSuperNode string, port string) {
	alpha := 3
	me := NewContact(NewRandomKademliaID(), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}, alpha, make(chan Message)}

	NodeNetwork.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode+port))
	NodeNetwork.SendFindContactMessage(&me)

}
