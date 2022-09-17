package d7024e


import (
	// "fmt"
	// "time"
)

// global var?
var NodeNetwork Network


func InitalizeSuperNode(id string, ip string){
	me := NewContact(NewKademliaID(id), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}
}

func InitalizeNode(ip string, idSuperNode string, ipSuperNode string, port string){
	me := NewContact(NewRandomKademliaID(), ip)
	NodeNetwork = Network{Kademlia{NewRoutingTable(me), 4, make(map[string][]byte)}}

	NodeNetwork.Kademlia.RoutingTable.AddContact(NewContact(NewKademliaID(idSuperNode), ipSuperNode+port))
	NodeNetwork.SendFindContactMessage(&me)

}

// func Pinger(n Network, me Contact) {

// 	for i := 0; i < 3; i++ {
// 		fmt.Println("Sending a ping ... NOW!")
// 		n.SendPingMessage(&me)
// 		time.Sleep(3 * time.Second)
// 	}

// }