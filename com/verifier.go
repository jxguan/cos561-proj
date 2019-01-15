package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"

	"github.com/jxguan/cos561-proj/lib"
	"gonum.org/v1/gonum/mat"
)

func Verify(sigma, phi, gamma, s mat.Dense, com, bit int) bool {
	prod := lib.BitMul(&sigma, &s)
	if !mat.Equal(&prod, &phi) {
		return false
	}
	prod = lib.BitMul(&gamma, &s)
	return com == ((int(prod.At(0, 0)) + bit) % 2)
}

func main() {

	var n, m, lambda int

	flag.IntVar(&n, "n", 256, "Security Parameter $n$")
	flag.IntVar(&m, "m", 512, "Security Parameter $m$")
	flag.IntVar(&lambda, "lambda", 64, "Security Parameter $\\lambda$")
	flag.Parse()

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	r := lib.CreateRecorder(conn, n, m, lambda, true)

	sigma, phi := r.EncryptZero()
	reader := bufio.NewReader(conn)

	// Receives the commitment
	var gamma mat.Dense
	_, err := gamma.UnmarshalBinaryFrom(reader)
	if err != nil {
		panic(err)
	}
	byteRead, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	com := int(byteRead)
	fmt.Println("Commitment Received")

	// Receives the opening
	var s mat.Dense
	_, err = s.UnmarshalBinaryFrom(reader)
	if err != nil {
		panic(err)
	}
	byteRead, err = reader.ReadByte()
	if err != nil {
		panic(err)
	}
	bit := int(byteRead)
	fmt.Printf("Opening to bit %d received. Verifying...\n", bit)

	if Verify(*sigma, *phi, gamma, s, com, bit) {
		fmt.Printf("Verified: Opening to %d is valid\n", bit)
	} else {
		fmt.Printf("Error! Opening to %d is not valid!\n", bit)
	}

	if !Verify(*sigma, *phi, gamma, s, com, 1-bit) {
		fmt.Printf("Verified: Opening to %d is not valid\n", 1-bit)
	} else {
		fmt.Printf("Error! Opening to %d is also valid!\n", 1-bit)
	}
}
