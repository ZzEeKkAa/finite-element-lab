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
	Lx, Rx float64
	Ly, Ry float64
	nx, ny int
	//gcord   []float64
	bcdof   []int
	bcval   []float64
}

func (fe *FiniteElement) Init(a, b, c float64, Lx, Rx, Ly, Ry float64, nx, ny int, ulx, urx, uly, ury func(x float64) float64) {
	//       uly
	//     +-----+
	// ulx |     | urx
	//     +-----+
	//       ury
	fe.a, fe.b, fe.c = a, b, c

	fe.nel = nx*ny // number of elements
	fe.nnel = 3 // number of nodes per element
	fe.nnode = (nx+1)*(ny+1) // number of nodes
	fe.ndof = 1

	// boundary conditions
	for i:=0; i<=nx; i++{
		// y=Ly
		fe.bcdof = append(fe.bcdof,i)
		fe.bcval = append(fe.bcval, uly(Lx+(Rx-Lx)*float64(i)/float64(nx)))
		// y=Ry
		fe.bcdof = append(fe.bcdof,i+(nx+1)*ny)
		fe.bcval = append(fe.bcval, ury(Lx+(Rx-Lx)*float64(i)/float64(nx)))
	}

	for i:=0; i<=ny; i++ {
		// x=Lx
		fe.bcdof = append(fe.bcdof, i*(nx+1))
		fe.bcval = append(fe.bcval, ulx(Ly+(Ry-Ly)*float64(i)/float64(ny)))
		// x=Rx
		fe.bcdof = append(fe.bcdof, i*(nx+1)+nx)
		fe.bcval = append(fe.bcval, urx(Ly+(Ry-Ly)*float64(i)/float64(ny)))
	}
}

func (fe *FiniteElement) gcord(i int)(float64, float64){
	return fe.Lx + (fe.Rx-fe.Lx)*float64(i%(fe.nx+1))/float64(fe.nx+1), fe.Ly + (fe.Ry-fe.Ly)*float64(i/(fe.nx+1))/float64(fe.ny+1)
}

func (fe *FiniteElement) nodes(a int) (int, int, int, int) {
	x:=a%fe.nx
	y:=a/fe.nx

	return y*(fe.nx+1)+x, y*(fe.nx+1)+x+1,(y+1)*(fe.nx+1)+x, (y+1)*(fe.nx+1)+x+1
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
