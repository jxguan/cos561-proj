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

	var n, m, lambda, size int

	flag.IntVar(&n, "n", 256, "Security Parameter $n$")
	flag.IntVar(&m, "m", 512, "Security Parameter $m$")
	flag.IntVar(&size, "size", 32, "Length of the bit string to be obliviously transferred")
	flag.Parse()
	lambda = 2 * size

	fmt.Println("Sender's Secret Strings:")
	x := make([][]int, 2)
	for i := 0; i < 2; i++ {
		fmt.Printf("\t x%d: ", i)
		x[i] = make([]int, size)
		for j := 0; j < size; j++ {
			x[i][j] = lib.RandomBit()
			fmt.Printf("%d", x[i][j])
		}
		fmt.Println()
	}

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
	c := int(byteRead)
	fmt.Println("Retrieval Query Received")

	for i := 0; i < size; i++ {
		for j := 0; j < 2; j++ {
			sigmaRow := mat.DenseCopyOf(sigma.Slice(2*i+j, 2*i+j+1, 0, m))
			seed := sigmaRow
			if x[j][i] == 1 {
				sum := lib.BitAdd(seed, &gamma)
				seed = &sum
			}
			_, err = seed.MarshalBinaryTo(conn)
			if err != nil {
				panic(err)
			}
			bit := (int(phi.At(2*i+j, 0)) + x[j][i]*(1-j+c)) % 2
			conn.Write([]byte{byte(bit)})
		}
	}

}
