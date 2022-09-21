package main

import (
	"fmt"
	"os"
	"time"
	."d7024e"
)

func init(){
	fmt.Println("Hello world")
	arg := os.Args

	idSuperNode := "0000000000000000000000000000000000000000"
	ipSuperNode := "172.20.0.2"
	port := ":80"
	ip := arg[1]
	ipAndPort := ip+port
	fmt.Println(ip)

	if ip == ipSuperNode{
		InitalizeSuperNode(idSuperNode, ipAndPort)
		NodeNetwork.Listen(ip, 80)	// keep this in main thread
	}else{
		InitalizeNode(ipAndPort, idSuperNode, ipSuperNode, port)
		contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
		
		Pinger(contact)
		// NodeNetwork.Listen(ip, 80)
		
	}

}

func main() {

	fmt.Println("")



}



func Pinger(me Contact) {

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		NodeNetwork.SendPingMessage(&me)
		time.Sleep(3 * time.Second)
	}

}



// test function to see if the super node is added as a contact
// retContact := network.Kademlia.LookupContact(&contact)
// if retContact[0].Address == contact.Address {
// 	fmt.Println("Jag existerar!")
// } else {
// 	fmt.Println(":(((((")
// }