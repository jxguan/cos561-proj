package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/mat"
)

func mod2(i int, j int, v float64) float64 {
	return float64(int(v) % 2)
}

func BitMul(a, b mat.Matrix) mat.Dense {
	var m mat.Dense
	m.Mul(a, b)
	m.Apply(mod2, &m)
	return m
}

func BitAdd(a, b mat.Matrix) mat.Dense {
	var m mat.Dense
	m.Add(a, b)
	m.Apply(mod2, &m)
	return m
}

func RandomBit() int {
	n, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

func PrintMatrix(m mat.Dense) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%d ", int(m.At(i, j)))
		}
		fmt.Println()
	}
}
