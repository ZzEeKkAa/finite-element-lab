package main

import (
	"fmt"
	"math"
)

func main() {
	var fe FiniteElement
	fe.Init(1, -3, 2, 5)

	sol := fe.Solve()
	esol := esol(5)

	for i, x := range sol {
		fmt.Printf("%7.4f %7.4f\n", x, esol[i])
	}
}

func esol(elements int) []float64 {
	esol := make([]float64, elements+1)

	c1, c2 := 0.5/math.Exp(1), -0.5*(1+1/math.Exp(1))

	x := 0.
	for i := range esol {
		esol[i] = c1*math.Exp(2*x) + c2*math.Exp(x) + 0.5
		x += 1 / float64(elements)
	}

	return esol
}
