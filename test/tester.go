// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"net"
// 	"time"
// )

// func main() {

// 	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 	rand.Seed(time.Now().Unix())

// 	num := 0
// 	for i := 0; i < 1; i++ {
// 		num++
// 		fmt.Println(num)
// 		con, _ := net.Dial("udp", "127.0.0.1:2000")
// 		buf := []byte("bla bla bla I am the packet")
// 		// _, err := con.Write(buf)

// 		for i := 0; i < len(buf); i++ {
// 			buf[i] = letters[rand.Intn(len(letters))]
// 		}

// 		_, err := con.Write(buf)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func main() {
	// Initialize the encoder and decoder.  Normally enc and dec would be
	// bound to network connections and the encoder and decoder would
	// run in different processes.
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	// HERE ARE YOUR BYTES!!!!
	fmt.Println(network.Bytes())

	// Decode (receive) the value.
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d,%d}\n", q.Name, *q.X, *q.Y)
}
