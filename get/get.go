package main

import (
	. "d7024e"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	arg := os.Args
	hash := arg[1]

	fmt.Println(hash) // remove later

	// Add code to check if hash is in correct format

	if len(hash) == 20 && isHexString(hash) {

		NodeNetwork.SendFindDataMessage(hash)

	} else {
		fmt.Println("Wrong hash format, please enter a valid hash ")
	}
}

func isHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}
