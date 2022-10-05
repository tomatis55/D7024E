package d7024e

import (
	"fmt"
)

var NodeNetwork Network

func InitalizeSuperNode(id string, ip string) {
	alpha := 3
	k := 4
	me := NewContact(NewKademliaID(id), ip)
	me.CalcDistance(me.ID)
	fmt.Println("Node ID: ", me.ID)
	NodeNetwork = NewNetwork(NewKademlia(NewRoutingTable(me), k), alpha)

	go NodeNetwork.Listen(ip, 80)
}

func InitalizeNode(ip string, idSuperNode string, ipSuperNode string, port string) {
	alpha := 3
	k := 4
	me := NewContact(NewRandomKademliaID(), ip)
	me.CalcDistance(me.ID)
	fmt.Println("Node ID: ", me.ID)
	NodeNetwork = NewNetwork(NewKademlia(NewRoutingTable(me), k), alpha)

	go NodeNetwork.Listen(ip, 80)
	superNode := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
	superNode.CalcDistance(me.ID)
	NodeNetwork.Kademlia.RoutingTable.AddContact(superNode)

	shortList := NodeNetwork.SendFindContactMessage(&me)
	if shortList.Len() == 0 {
		NodeNetwork.SendFindContactMessage(&me)
	}

}
