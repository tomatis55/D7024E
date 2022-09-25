package main

var NodeNetwork Network

func InitalizeSuperNode(id string, ip string) {
	alpha := 3
	k := 4
	me := NewContact(NewKademliaID(id), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), k, make(map[string][]byte)}, alpha, make(chan Message, alpha)}

	go NodeNetwork.Listen(ip, 80)
}

func InitalizeNode(ip string, idSuperNode string, ipSuperNode string, port string) {
	alpha := 3
	k := 4
	me := NewContact(NewRandomKademliaID(), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), k, make(map[string][]byte)}, alpha, make(chan Message, alpha)}

	go NodeNetwork.Listen(ip, 80)
	superNode := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
	superNode.CalcDistance(me.ID)
	NodeNetwork.Kademlia.RoutingTable.AddContact(superNode)

	NodeNetwork.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode+port))
	NodeNetwork.SendFindContactMessage(&me)

}
