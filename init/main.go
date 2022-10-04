package main

import (
	"bufio"
	. "d7024e"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	arg := os.Args

	idSuperNode := "0000000000000000000000000000000000000001"
	ipSuperNode := "172.20.0.2"
	port := ":80"
	ip := arg[1]
	ipAndPort := ip + port
	fmt.Println(ip)

	if ip == ipSuperNode {
		InitalizeSuperNode(idSuperNode, ipAndPort)

	} else {
		if len(os.Args) > 2 { // if another super node was specified
			InitalizeNode(ipAndPort, arg[2], arg[3], port) // arg[2] = id of node to connect to, arg[3] = ip of node to connect to
		} else { // if no node was specified, use the standard super node
			InitalizeNode(ipAndPort, idSuperNode, ipSuperNode, port)
		}
	}
}

func main() {

	r := bufio.NewReader(os.Stdin)

	loop := true
	for loop {

		input, _, _ := r.ReadLine()
		inputSlices := strings.SplitN(string(input), " ", 2)

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

		case inputSlices[0] == "info":
			Info()

		case inputSlices[0] == "superinfo":
			SuperInfo()

		case inputSlices[0] == "find" && len(inputSlices) == 2:
			FindNode(inputSlices[1])

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
