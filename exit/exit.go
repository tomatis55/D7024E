package main

import (
	. "d7024e"
	"fmt"
)

func main() {

	fmt.Println("in exit now") // remove later

	NodeNetwork.SendTerminateNodeMessage()
}
