package lib

import (
	"net"

	"gonum.org/v1/gonum/mat"
)

// Streamer stores the essential information needed for a streamer to execute an
// EncryptZero protocol.
type Keeper struct {
	conn       net.Conn // The already opened TCP connection
	k          *mat.Dense
	s          *mat.Dense
	n          int
	m          int
	isEnhanced bool
}

func CreateKeeper(conn net.Conn, n int, m int, isEnhanced bool) *Keeper {
	kp := new(Keeper)
	kp.conn = conn
	kp.isEnhanced = isEnhanced
	kp.n = n
	kp.m = m
	randomBits := make([]float64, n)
	for i := range randomBits {
		randomBits[i] = float64(RandomBit())
	}
	kp.k = mat.NewDense(n, 1, randomBits)
	if isEnhanced {
		randomBits := make([]float64, m)
		for i := range randomBits {
			randomBits[i] = float64(RandomBit())
		}
		kp.s = mat.NewDense(m, 1, randomBits)
	}
	return kp
}

func (kp *Keeper) EncryptZero() *mat.Dense {
	for i := 0; i < kp.m; i++ {
		randomBits := make([]float64, kp.n)
		for i := range randomBits {
			randomBits[i] = float64(RandomBit())
		}
		r := mat.NewDense(1, kp.n, randomBits)

		result := BitMul(r, kp.k)
		a := int(result.At(0, 0))
		if kp.isEnhanced {
			a = (a + int(kp.s.At(i, 0))) % 2
		}
		_, err := r.MarshalBinaryTo(kp.conn)
		if err != nil {
			panic(err)
		}
		kp.conn.Write([]byte{byte(a)})
	}
	if kp.isEnhanced {
		_, err := kp.k.MarshalBinaryTo(kp.conn)
		if err != nil {
			panic(err)
		}
		return kp.s
	}
	return kp.k
}
