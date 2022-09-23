package main

import (
	"fmt"
	"os"
	"time"
	"bufio"
	."d7024e"
	"strings"
	"os/exec"
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

	// }else if (ip == "172.20.0.3"){
	// 	InitalizeSuperNode("0000000000000000000000000000000000000001", "172.20.0.3:80")

	}else{
		if len(os.Args) > 2{	// if another super node was specified
			InitalizeNode(ipAndPort, arg[2], arg[3], port)		// arg[2] = id of node to connect to, arg[3] = ip of node to connect to
		}else{		// if no node was specified, use the standard super node
			InitalizeNode(ipAndPort, idSuperNode, ipSuperNode, port)
		}
		//contact := NewContact(NewKademliaID(idSuperNode), ipSuperNode+port)
		//Pinger(contact)

	}
}

func main() {
	arg := os.Args
	ip := arg[1]

	contact := NodeNetwork.Kademlia.RoutingTable.FindClosestContacts(NewKademliaID("0000000000000000000000000000000000000001"), 2)
	fmt.Println("Number of contacts: ", len(contact))
	if len(contact) > 0{
		fmt.Println("Closest contact: ",contact[0].Address)
	}
	

	go NodeNetwork.Listen(ip, 80)

	loop := true
	for loop{
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
			loop = false
		
		case inputSlices[0] == "ping" && len(inputSlices) == 2:
			Ping(inputSlices[1])

		default:
			fmt.Println("Not a valid command")
		}
	}

	exec.Command("kill -s SIGTERM 1")

}



func Pinger(contact Contact) {

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		NodeNetwork.SendPingMessage(&contact)
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