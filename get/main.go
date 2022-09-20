package main

import (
	. "d7024e"
	"fmt"
	"os"
)

func main() {
	arg := os.Args

	id := NewRandomKademliaID()
	_ = id
	//fmt.Println("Generated id: ", id)

	fmt.Println(arg[1])
}
