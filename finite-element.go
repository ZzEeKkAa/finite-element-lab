package main

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

type FiniteElement struct {
	nodes    []Node
	elements []Element
	boundary []NodeVal
	nel      int
	nnel     int
	ndof     int
	nnode    int
	sdof     int
}

func (fe *FiniteElement) Init(nodes []Node, elements []Element, boundary []NodeVal) {
	fe.nodes = nodes
	fe.elements = elements
	fe.boundary = boundary
	fe.nel = len(elements)
	fe.nnode = len(nodes)
	fe.nnel = 4
	fe.ndof = 1
	fe.sdof = fe.nnode * fe.ndof
}

func (fe *FiniteElement) Solve() []float64 {
	var (
		ff = mat64.NewDense(fe.sdof, 1, nil)
		kk = mat64.NewDense(fe.sdof, fe.sdof, nil)
	)

	for i := 0; i < fe.nel; i++ {
		element := fe.elements[i]
		nodes := element.Nodes(fe.nodes)

		xBar := mat64.NewDense(4, 4, []float64{
			1, nodes[0].X, nodes[0].Y, nodes[0].Z,
			1, nodes[1].X, nodes[1].Y, nodes[1].Z,
			1, nodes[2].X, nodes[2].Y, nodes[2].Z,
			1, nodes[3].X, nodes[3].Y, nodes[3].Z,
		})

		var xInv = &mat64.Dense{}
		xInv.Inverse(xBar)

		vol := mat64.Det(xBar) / 6.
		if vol < 0 {
			vol = -vol
		}

		//fmt.Println(i, vol)
		//
		//fmt.Println(element)

		k11 := vol * (xInv.At(1, 0)*xInv.At(1, 0) + xInv.At(2, 0)*xInv.At(2, 0) + xInv.At(3, 0)*xInv.At(3, 0))
		k12 := vol * (xInv.At(1, 0)*xInv.At(1, 1) + xInv.At(2, 0)*xInv.At(2, 1) + xInv.At(3, 0)*xInv.At(3, 1))
		k13 := vol * (xInv.At(1, 0)*xInv.At(1, 2) + xInv.At(2, 0)*xInv.At(2, 2) + xInv.At(3, 0)*xInv.At(3, 2))
		k14 := vol * (xInv.At(1, 0)*xInv.At(1, 3) + xInv.At(2, 0)*xInv.At(2, 3) + xInv.At(3, 0)*xInv.At(3, 3))

		k22 := vol * (xInv.At(1, 1)*xInv.At(1, 1) + xInv.At(2, 1)*xInv.At(2, 1) + xInv.At(3, 1)*xInv.At(3, 1))
		k23 := vol * (xInv.At(1, 1)*xInv.At(1, 2) + xInv.At(2, 1)*xInv.At(2, 2) + xInv.At(3, 1)*xInv.At(3, 2))
		k24 := vol * (xInv.At(1, 1)*xInv.At(1, 3) + xInv.At(2, 1)*xInv.At(2, 3) + xInv.At(3, 1)*xInv.At(3, 3))

		k33 := vol * (xInv.At(1, 2)*xInv.At(1, 2) + xInv.At(2, 2)*xInv.At(2, 2) + xInv.At(3, 2)*xInv.At(3, 2))
		k34 := vol * (xInv.At(1, 2)*xInv.At(1, 3) + xInv.At(2, 2)*xInv.At(2, 3) + xInv.At(3, 2)*xInv.At(3, 3))

		k44 := vol * (xInv.At(1, 3)*xInv.At(1, 3) + xInv.At(2, 3)*xInv.At(2, 3) + xInv.At(3, 3)*xInv.At(3, 3))

		kk.Set(element.N1, element.N1, kk.At(element.N1, element.N1)+k11)
		kk.Set(element.N1, element.N2, kk.At(element.N1, element.N2)+k12)
		kk.Set(element.N1, element.N3, kk.At(element.N1, element.N3)+k13)
		kk.Set(element.N1, element.N4, kk.At(element.N1, element.N4)+k14)
		kk.Set(element.N2, element.N1, kk.At(element.N2, element.N1)+k12)
		kk.Set(element.N2, element.N2, kk.At(element.N2, element.N2)+k22)
		kk.Set(element.N2, element.N3, kk.At(element.N2, element.N3)+k23)
		kk.Set(element.N2, element.N4, kk.At(element.N2, element.N4)+k24)
		kk.Set(element.N3, element.N1, kk.At(element.N3, element.N1)+k13)
		kk.Set(element.N3, element.N2, kk.At(element.N3, element.N2)+k23)
		kk.Set(element.N3, element.N3, kk.At(element.N3, element.N3)+k33)
		kk.Set(element.N3, element.N4, kk.At(element.N3, element.N4)+k34)
		kk.Set(element.N4, element.N1, kk.At(element.N4, element.N1)+k14)
		kk.Set(element.N4, element.N2, kk.At(element.N4, element.N2)+k24)
		kk.Set(element.N4, element.N3, kk.At(element.N4, element.N3)+k34)
		kk.Set(element.N4, element.N4, kk.At(element.N4, element.N4)+k44)
	}

	//printDense(kk)
	fe.feaplyc2(kk, ff)

	//fmt.Println()

	//printDense(kk)
	//printDense(ff)

	u := &mat64.Dense{}
	u.Solve(kk, ff)

	return u.RawMatrix().Data
}

func (fe *FiniteElement) feode2l(h float64, k [][]float64) {
	//a1, a2, a3 := -fe.a/h, fe.b/2, fe.c*h/6
	//k[0][0] += a1 - a2 + 2*a3
	//k[0][1] += -a1 + a2 + a3
	//k[1][0] += -a1 - a2 + a3
	//k[1][1] += a1 + a2 + 2*a3
}

func (fe *FiniteElement) fef1l(h float64, f []float64) {
	f[0] += h / 2
	f[1] += h / 2
}

func (fe *FiniteElement) feaplyc2(kk *mat64.Dense, ff *mat64.Dense) {
	for _, boubdary := range fe.boundary {
		for j := 0; j < fe.sdof; j++ {
			kk.Set(boubdary.node, j, 0)
		}
		kk.Set(boubdary.node, boubdary.node, 1)
		//fmt.Println(boubdary.node, boubdary.val)
		ff.Set(boubdary.node, 0, boubdary.val)
	}
}

func printDense(dense *mat64.Dense) {
	r, c := dense.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%6.2f", dense.At(i, j))
		}

		fmt.Println()
	}
}
