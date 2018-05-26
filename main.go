package main

import (
	"fmt"
	"math"
)

func main() {


	for elements:= 1; elements<5000; elements+=1000 {

		var (
			fe FiniteElement
			//elements int = 15
		)

		fe.Init(1, -3, 2, elements)
		esol := esol(elements)
		sol := fe.Solve()

		//fmt.Println(sol)
		//fmt.Println(esol)

		var totalError float64
		for i, x := range sol {
			err := x - esol[i]
			totalError += err * err
			//fmt.Printf("%7.4f %7.4f\n", x, esol[i])
		}

		totalError = math.Sqrt(totalError / float64(elements+1))

		fmt.Println(totalError)
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
