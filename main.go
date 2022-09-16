package main

import "time"

func main() {
	go Listen("127.0.0.1", 2000)
	Pinger()
}

func Pinger() {
	id := NewKademliaID("5465747261687964726F63616E6E6162696E6F6C")
	d := id.CalcDistance(id)
	me := Contact{ID: id, Address: "127.0.0.1:2000", distance: d}
	kad := Kademlia{NewRoutingTable(me)}

	n := Network{kad}

	for i := 0; i < 3; i++ {
		n.SendPingMessage(&me)
		time.Sleep(3 * time.Second)
	}

}
