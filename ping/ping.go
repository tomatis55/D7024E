package main

import (
	. "d7024e"
	"fmt"
	"time"
	"os"
)

func main() {

	arg := os.Args
	ip := arg[1]+":80"
	fmt.Println(ip)

	contact := NewContact(NewRandomKademliaID(), ip)

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		NodeNetwork.SendPingMessage(&contact)
		time.Sleep(3 * time.Second)
	}
}