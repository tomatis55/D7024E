package main

import (
	"fmt"
	"os"
	"time"
	"bufio"
	."d7024e"
	"strings"
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

		// contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "172.20.0.3:80")
		// Pinger(contact)

	}else{
		
		if len(os.Args) > 2{
			InitalizeNode(ipAndPort, arg[2], arg[3], port)
		}else{
			InitalizeNode(ipAndPort, idSuperNode, ipSuperNode, port)
		}
		//contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
		//Pinger(contact)

	}
}

func main() {
	arg := os.Args
	ip := arg[1]

	go NodeNetwork.Listen(ip, 80)

	for{
		r := bufio.NewReader(os.Stdin)
		
		input, _, _ := r.ReadLine()

		inputSlices := strings.Split(string(input), " ")

		switch {
		case inputSlices[0] == "get" && len(inputSlices) == 2:
			Get(inputSlices[1])
			
		case inputSlices[0] == "put" && len(inputSlices) == 2:
			Put(inputSlices[1])
		
		case inputSlices[0] == "exit":
			Exit()
		
		case inputSlices[0] == "ping" && len(inputSlices) == 2:
			Ping(inputSlices[1])

		default:
			fmt.Println("Not a valid command")
		}
	}



}



func Pinger(me Contact) {

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		NodeNetwork.SendPingMessage(&me)
		time.Sleep(30 * time.Second)
	}

}



// test function to see if the super node is added as a contact
// retContact := network.Kademlia.LookupContact(&contact)
// if retContact[0].Address == contact.Address {
// 	fmt.Println("Jag existerar!")
// } else {
// 	fmt.Println(":(((((")
// }