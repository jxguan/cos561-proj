package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/jxguan/cos561-proj/lib"
)

func main() {

	var n, m, lambda int

	flag.IntVar(&n, "n", 256, "Security Parameter $n$")
	flag.IntVar(&m, "m", 512, "Security Parameter $m$")
	flag.IntVar(&lambda, "keylen", 64, "Length of Derived Key")
	flag.Parse()

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	r := lib.CreateRecorder(conn, n, m, lambda, false)

	psi, kappa := r.EncryptZero()

	_, err := psi.MarshalBinaryTo(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Derived Key:")
	len, _ := kappa.Dims()
	for i := 0; i < len; i++ {
		fmt.Printf("%d", int(kappa.At(i, 0)))
	}
	fmt.Println()
}
