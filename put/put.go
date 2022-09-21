package main

import (
	"fmt"
	."d7024e"
	"os"
)

func main() {
	arg := os.Args
	dataStr := arg[1]

	// Add code to check if dataStr is in correct format

	if len(dataStr) <= 255{
		data := []byte(dataStr)	

		NodeNetwork.SendStoreMessage(data)

	}else{
		fmt.Println("Too large data string")
	}


}