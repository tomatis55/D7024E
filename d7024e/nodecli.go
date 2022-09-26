package d7024e

import (
	"encoding/hex"
	"fmt"
	"time"
)

func Get(hash string) {

	if len(hash) == 40 && isHexString(hash) {

		NodeNetwork.SendFindDataMessage(hash)

	} else {
		fmt.Println("Wrong hash format, please enter a valid hash ")
	}
}

func isHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func Put(dataStr string) {

	if len(dataStr) <= 255 {
		data := []byte(dataStr)

		NodeNetwork.SendStoreMessage(data)

	} else {
		fmt.Println("Too large data string")
	}
}

func Exit() {
	NodeNetwork.SendTerminateNodeMessage()
}

func Ping(ip string) {
	ip = ip + ":80"
	fmt.Println(ip)

	contact := NewContact(NewRandomKademliaID(), ip)

	for i := 0; i < 3; i++ {
		fmt.Println("Sending a ping ... NOW!")
		NodeNetwork.SendPingMessage(&contact)
		time.Sleep(3 * time.Second)
	}
}

func Info() {
	ip := NodeNetwork.Kademlia.RoutingTable.me.Address
	id := NodeNetwork.Kademlia.RoutingTable.me.ID
	fmt.Println("Node IP: ", ip)
	fmt.Println("Node ID: ", id)
}

func SuperInfo() {
	ip := NodeNetwork.Kademlia.RoutingTable.me.Address
	id := NodeNetwork.Kademlia.RoutingTable.me.ID
	fmt.Println("Node IP: ", ip)
	fmt.Println("Node ID: ", id)
	fmt.Println("======== All known contacts ========")
	for _, x := range NodeNetwork.Kademlia.GetAllContacts().contacts {
		fmt.Println(x.String())
	}
}
