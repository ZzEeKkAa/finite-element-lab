package main

import (
	"github.com/gonum/matrix/mat64"
)

type FiniteElement struct {
	a, b, c float64
	el      []float64
	nel     int
	nnel    int
	ndof    int
	nnode   int
	sdof    int
	gcord   []float64
	bcdof   []int
	bcval   []float64
}

func (fe *FiniteElement) Init(a, b, c float64, elements int) {
	fe.a, fe.b, fe.c = a, b, c

	fe.nel = elements
	fe.nnel = 2
	fe.nnode = elements + 1
	fe.ndof = 1
	for i, d := 0.0, 1/float64(elements); i < 1; i += d {
		fe.gcord = append(fe.gcord, i)
	}
	if len(fe.gcord) < elements+1 {
		fe.gcord = append(fe.gcord, 1)
	}
	fe.gcord[elements] = 1

	//boundary conditions
	fe.bcdof = []int{0, elements}
	fe.bcval = []float64{0, 0}
}

func (fe *FiniteElement) nodes(a int) (int, int) {
	return a, a + 1
}

func (fe *FiniteElement) Solve() []float64 {
	var ff []float64
	var kkData []float64
	var kk [][]float64

	ff = make([]float64, fe.nnode)
	kkData = make([]float64, fe.nnode*fe.nnode)
	kk = make([][]float64, fe.nnode)

	for i := range kk {
		kk[i] = kkData[i*fe.nnode : (i+1)*fe.nnode]
	}

	for i := 0; i < fe.nel; i++ {
		nl, nr := fe.nodes(i)
		xl, xr := fe.gcord[nl], fe.gcord[nr]
		h := xr - xl
		k := make([][]float64, 2)
		k[0] = kk[i][i : i+2]
		k[1] = kk[i+1][i : i+2]
		f := ff[i : i+2]
		fe.feode2l(h, k)
		fe.fef1l(h, f)
	}

	fe.feaplyc2(kk, ff)

	A := mat64.NewDense(fe.nnode, fe.nnode, kkData)
	b := mat64.NewDense(fe.nnode, 1, ff)
	var x = &mat64.Dense{}
	x.Solve(A, b)

	return x.RawMatrix().Data
}

func (fe *FiniteElement) feode2l(h float64, k [][]float64) {
	a1, a2, a3 := -fe.a/h, fe.b/2, fe.c*h/6
	k[0][0] += a1 - a2 + 2*a3
	k[0][1] += -a1 + a2 + a3
	k[1][0] += -a1 - a2 + a3
	k[1][1] += a1 + a2 + 2*a3
}

func (fe *FiniteElement) fef1l(h float64, f []float64) {
	f[0] += h / 2
	f[1] += h / 2
}

func (fe *FiniteElement) feaplyc2(kk [][]float64, ff []float64) {
	for i := range fe.bcdof {
		c := fe.bcdof[i]
		for j := range kk[c] {
			kk[c][j] = 0
		}

		kk[c][c] = 1
		ff[c] = fe.bcval[i]
	}
}
