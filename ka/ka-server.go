package main

import (
	"bufio"
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
	flag.IntVar(&lambda, "keylen", 64, "Length of Derived Key")
	flag.Parse()
	fmt.Println("Waiting for client to connect...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	k := lib.CreateKeeper(conn, n, m, false)

	s := k.EncryptZero()

	var psi mat.Dense
	_, err := psi.UnmarshalBinaryFrom(bufio.NewReader(conn))
	if err != nil {
		panic(err)
	}

	kappa := lib.BitMul(&psi, s)

	fmt.Println("Derived Key:")
	len, _ := kappa.Dims()
	for i := 0; i < len; i++ {
		fmt.Printf("%d", int(kappa.At(i, 0)))
	}
	fmt.Println()
}
