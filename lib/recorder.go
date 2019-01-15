package lib

import (
	"bufio"
	"net"

	"gonum.org/v1/gonum/mat"
)

// Recorder stores the essential information needed for a recorder to execute an
// EncryptZero protocol.
type Recorder struct {
	conn       net.Conn // The already opened TCP connection
	n          int
	lambda     int
	sigma      *mat.Dense
	psi        *mat.Dense
	kappa      *mat.Dense
	m          int
	isEnhanced bool
}

func CreateRecorder(conn net.Conn, n, m, lambda int, isEnhanced bool) *Recorder {
	r := new(Recorder)
	r.conn = conn
	r.isEnhanced = isEnhanced
	r.n = n
	r.m = m
	r.lambda = lambda
	r.psi = mat.NewDense(lambda, n, nil)
	r.kappa = mat.NewDense(lambda, 1, nil)
	if isEnhanced {
		r.sigma = mat.NewDense(lambda, m, nil)
	}
	return r
}

func (rc *Recorder) EncryptZero() (*mat.Dense, *mat.Dense) {
	reader := bufio.NewReader(rc.conn)
	for i := 0; i < rc.m; i++ {
		var r mat.Dense
		_, err := r.UnmarshalBinaryFrom(reader)
		if err != nil {
			panic(err)
		}
		randomBits := make([]float64, rc.lambda)
		for i := range randomBits {
			randomBits[i] = float64(RandomBit())
		}
		sigmaCol := mat.NewDense(rc.lambda, 1, randomBits)
		byteRead, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		a := int(byteRead)
		prod := BitMul(sigmaCol, &r)
		sum := BitAdd(&prod, rc.psi)
		rc.psi = mat.DenseCopyOf(&sum)
		if a == 1 {
			sum := BitAdd(sigmaCol, rc.kappa)
			rc.kappa = mat.DenseCopyOf(&sum)
		}
		if rc.isEnhanced {
			rc.sigma.SetCol(i, randomBits)
		}
	}
	if rc.isEnhanced {
		var k mat.Dense
		_, err := k.UnmarshalBinaryFrom(reader)
		if err != nil {
			panic(err)
		}
		prod := BitMul(rc.psi, &k)
		sum := BitAdd(rc.kappa, &prod)
		phi := mat.DenseCopyOf(&sum)
		return rc.sigma, phi
	}
	return rc.psi, rc.kappa
}
