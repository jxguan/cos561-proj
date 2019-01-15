package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/jxguan/cos561-proj/lib"
	"gonum.org/v1/gonum/mat"
)

func main() {

	var n, m, lambda int

	flag.IntVar(&n, "n", 256, "Security Parameter $n$")
	flag.IntVar(&m, "m", 512, "Security Parameter $m$")
	flag.IntVar(&lambda, "lambda", 64, "Security Parameter $\\lambda$")
	flag.Parse()
	fmt.Println("Waiting for client to connect...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	bit := lib.RandomBit()
	fmt.Printf("Secret Bit: %d\n", bit)

	k := lib.CreateKeeper(conn, n, m, true)

	s := k.EncryptZero()

	randomBits := make([]float64, m)
	for i := range randomBits {
		randomBits[i] = float64(lib.RandomBit())
	}
	gamma := mat.NewDense(1, m, randomBits)
	prod := lib.BitMul(gamma, s)
	com := (int(prod.At(0, 0)) + bit) % 2

	// Sends over the commitment
	_, err := gamma.MarshalBinaryTo(conn)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte{byte(com)})

	// Sends over the opening
	_, err = s.MarshalBinaryTo(conn)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte{byte(bit)})
}
