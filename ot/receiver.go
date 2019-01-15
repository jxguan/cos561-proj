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

	var n, m, size int

	flag.IntVar(&n, "n", 256, "Security Parameter $n$")
	flag.IntVar(&m, "m", 512, "Security Parameter $m$")
	flag.IntVar(&size, "size", 32, "Length of the bit string to be obliviously transferred")
	flag.Parse()

	fmt.Println("Waiting for sender to connect...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	bit := lib.RandomBit()
	fmt.Printf("String Requested: x%d\n", bit)

	k := lib.CreateKeeper(conn, n, m, true)

	s := k.EncryptZero()

	randomBits := make([]float64, m)
	for i := range randomBits {
		randomBits[i] = float64(lib.RandomBit())
	}
	gamma := mat.NewDense(1, m, randomBits)
	prod := lib.BitMul(gamma, s)
	com := (int(prod.At(0, 0)) + bit) % 2

	// Sends over the commitment of the string requested
	_, err := gamma.MarshalBinaryTo(conn)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte{byte(com)})

	// Prepare to store the string
	x := make([][]int, 2)
	for i := 0; i < 2; i++ {
		x[i] = make([]int, size)
	}

	reader := bufio.NewReader(conn)
	for i := 0; i < size; i++ {
		for j := 0; j < 2; j++ {
			var seed mat.Dense
			_, err := seed.UnmarshalBinaryFrom(reader)
			if err != nil {
				panic(err)
			}
			byteRead, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}
			c := int(byteRead)

			prod := lib.BitMul(&seed, s)
			x[j][i] = (int(prod.At(0, 0)) + c) % 2
		}
	}

	fmt.Println("Retrieved Strings:")
	for i := 0; i < 2; i++ {
		fmt.Printf("\t x%d: ", i)
		for j := 0; j < size; j++ {
			fmt.Printf("%d", x[i][j])
		}
		fmt.Println()
	}

}
