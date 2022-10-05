package d7024e

import (
	"fmt"
)

var NodeNetwork Network

func InitalizeSuperNode(id string, ip string, port int) *Network {
	alpha := 3
	k := 4
	me := NewContact(NewKademliaID(id), ip+":"+fmt.Sprint(port))
	me.CalcDistance(me.ID)
	fmt.Println("Node ID: ", me.ID)
	NodeNetwork = NewNetwork(NewKademlia(NewRoutingTable(me), k), alpha)

	go NodeNetwork.Listen(ip, port)
	return &NodeNetwork
}

func InitalizeNode(ip string, port int, idSuperNode string, ipSuperNode string, portSuperNode string) *Network {
	alpha := 3
	k := 4
	me := NewContact(NewRandomKademliaID(), ip+":"+fmt.Sprint(port))
	me.CalcDistance(me.ID)
	fmt.Println("Node ID: ", me.ID)
	NodeNetwork = NewNetwork(NewKademlia(NewRoutingTable(me), k), alpha)

	go NodeNetwork.Listen(ip, port)
	superNode := NewContact(NewKademliaID(idSuperNode), ipSuperNode+":"+portSuperNode)
	superNode.CalcDistance(me.ID)
	NodeNetwork.Kademlia.RoutingTable.AddContact(superNode)

	shortList := NodeNetwork.SendFindContactMessage(&me)
	if shortList.Len() == 0 {
		NodeNetwork.SendFindContactMessage(&me)
	}

	return &NodeNetwork
}
