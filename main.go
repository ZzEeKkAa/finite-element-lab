package main

import (
	"fmt"
	"math"
	"runtime"
	"sort"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

type Node struct {
	X, Y, Z float64
}

type NodeVal struct {
	node int
	val  float64
}

type Element struct {
	N1, N2, N3, N4 int
}

func (el Element) Nodes(nodes []Node) []Node {
	return []Node{nodes[el.N1], nodes[el.N2], nodes[el.N3], nodes[el.N4]}
}
func (el Element) NodeIndexes(nodes []Node) []int {
	return []int{el.N1, el.N2, el.N3, el.N4}
}

func main() {
	u := func(x, y, z float64) float64 {
		return x*x - 0.5*(y*y+z*z)
	}

	boundaryNodeIndexes := []int{0, 1, 2, 3, 4, 5}

	nodes := []Node{{0, 0, 0}, {1, 0, 0}, {0.5, 0.5, 0}, {0, 1, 0}, {1, 1, 0}, {0.5, 0.5, 1}, {0.5, 0.5, 0.5}}
	elements := []Element{
		{3, 0, 6, 5},
		{3, 0, 2, 6},
		{0, 1, 6, 5},
		{0, 1, 2, 6},
		{1, 4, 6, 5},
		{1, 4, 2, 6},
		{4, 3, 6, 5},
		{4, 3, 2, 6},
	}
	fmt.Println(len(boundaryNodeIndexes))
	//boundary := []NodeVal{{0, 0}, {1, 20}, {3, 50}, {4, 100}}
	//for i := 0; i < 2; i++ {
	//elements, nodes, boundaryNodeIndexes = split(elements, nodes, boundaryNodeIndexes)

	elements, nodes, boundaryNodeIndexes = split(elements, nodes, boundaryNodeIndexes)
	fmt.Println(len(boundaryNodeIndexes))
	elements, nodes, boundaryNodeIndexes = split2(elements, nodes, boundaryNodeIndexes)
	fmt.Println(len(boundaryNodeIndexes))
	//elements, nodes, boundaryNodeIndexes = split(elements, nodes, boundaryNodeIndexes)
	//fmt.Println(len(boundaryNodeIndexes))
	//elements, nodes, boundaryNodeIndexes = split(elements, nodes, boundaryNodeIndexes)
	//fmt.Println(len(boundaryNodeIndexes))
	//elements, nodes, boundaryNodeIndexes = split(elements, nodes, boundaryNodeIndexes)
	//fmt.Println(len(boundaryNodeIndexes))
	//}

	show3D(nodes, elements)
	//
	//return

	boundary := buildBoundary(nodes, boundaryNodeIndexes, u)

	fe := &FiniteElement{}

	fe.Init(nodes, elements, boundary)
	esol := esol(nodes, u)
	sol := fe.Solve()

	//fmt.Println(sol)
	//fmt.Println(esol)

	var totalError float64
	var minErr = math.MaxFloat64
	var maxErr = -math.MaxFloat64
	var num int
	for i, x := range sol {
		boundary := "-"

		if inArr(boundaryNodeIndexes, i) {
			boundary = "+"
		} else {
			num++
			err := math.Abs(x - esol[i])
			totalError += err * err
			if minErr > err {
				minErr = err
			}
			if maxErr < err {
				maxErr = err
			}
		}

		fmt.Printf("%2d %s %7.4f %7.4f\n", i, boundary, x, esol[i])
		//fmt.Printf("%d %7.4f\n", i, x)
	}

	totalError = math.Sqrt(totalError / float64(num))

	fmt.Println(totalError, minErr, maxErr)
	//}
}

func esol(nodes []Node, u func(x, y, z float64) float64) []float64 {
	esol := make([]float64, len(nodes))

	for i, node := range nodes {
		esol[i] = u(node.X, node.Y, node.Z)
	}

	return esol
}

func buildBoundary(nodes []Node, nodeIndexes []int, u func(x, y, z float64) float64) []NodeVal {
	var res []NodeVal
	for _, i := range nodeIndexes {
		res = append(res, NodeVal{i, u(nodes[i].X, nodes[i].Y, nodes[i].Z)})
	}

	return res
}

func split(elements []Element, nodes []Node, boundaryNodeIndexes []int) ([]Element, []Node, []int) {
	var newElements []Element
	for _, el := range elements {
		var node1 Node
		node1.X = (nodes[el.N1].X + nodes[el.N2].X + nodes[el.N3].X + nodes[el.N4].X) / 4.
		node1.Y = (nodes[el.N1].Y + nodes[el.N2].Y + nodes[el.N3].Y + nodes[el.N4].Y) / 4.
		node1.Z = (nodes[el.N1].Z + nodes[el.N2].Z + nodes[el.N3].Z + nodes[el.N4].Z) / 4.

		var b1, b2, b3, b4 bool

		if !inArr(boundaryNodeIndexes, el.N1) {
			b1 = true
			b2 = true
			b3 = true
		}
		if !inArr(boundaryNodeIndexes, el.N2) {
			b1 = true
			b2 = true
			b4 = true
		}
		if !inArr(boundaryNodeIndexes, el.N3) {
			b1 = true
			b3 = true
			b4 = true
		}
		if !inArr(boundaryNodeIndexes, el.N4) {
			b2 = true
			b3 = true
			b4 = true
		}
		var n1, n2, n3 int

		nodes = append(nodes, node1)
		nodeID := len(nodes) - 1
		if b1 {
			newElements = append(newElements, Element{el.N1, el.N2, el.N3, nodeID})
		} else {
			n1, n2, n3 = el.N1, el.N2, el.N3
		}
		if b2 {
			newElements = append(newElements, Element{el.N1, el.N2, nodeID, el.N4})
		} else {
			n1, n2, n3 = el.N1, el.N2, el.N4
		}
		if b3 {
			newElements = append(newElements, Element{el.N1, nodeID, el.N3, el.N4})
		} else {
			n1, n2, n3 = el.N1, el.N3, el.N4
		}
		if b4 {
			newElements = append(newElements, Element{nodeID, el.N2, el.N3, el.N4})
		} else {
			n1, n2, n3 = el.N2, el.N3, el.N4
		}

		if n1 != n2 {
			var node2 Node
			node2.X = (nodes[n1].X + nodes[n2].X + nodes[n3].X) / 3.
			node2.Y = (nodes[n1].Y + nodes[n2].Y + nodes[n3].Y) / 3.
			node2.Z = (nodes[n1].Z + nodes[n2].Z + nodes[n3].Z) / 3.

			node2ID := len(nodes)
			nodes = append(nodes, node2)

			boundaryNodeIndexes = append(boundaryNodeIndexes, node2ID)

			newElements = append(newElements,
				Element{nodeID, node2ID, n1, n2},
				Element{nodeID, node2ID, n2, n3},
				Element{nodeID, node2ID, n1, n3},
			)
		}
	}

	return newElements, nodes, boundaryNodeIndexes
}

func split2(elements []Element, nodes []Node, boundaryNodeIndexes []int) ([]Element, []Node, []int) {
	var newElements []Element
	for _, el := range elements {
		var extraCount int
		var extraNode int
		if inArr(boundaryNodeIndexes, el.N1) {
			extraCount++
			extraNode += el.N1
		}
		if inArr(boundaryNodeIndexes, el.N2) {
			extraCount++
			extraNode += el.N2
		}
		if inArr(boundaryNodeIndexes, el.N3) {
			extraCount++
			extraNode += el.N3
		}
		if inArr(boundaryNodeIndexes, el.N4) {
			extraCount++
			extraNode += el.N4
		}

		if extraCount == 3 {
			extraNode = el.N1 + el.N2 + el.N3 + el.N4 - extraNode

			switch extraNode {
			case el.N1:
				el.N1, el.N4 = el.N4, el.N1
			case el.N2:
				el.N2, el.N4 = el.N4, el.N2
			case el.N3:
				el.N3, el.N4 = el.N4, el.N3
			}
		}

		var (
			N5  = midNode(&nodes, el.N1, el.N2)
			N6  = midNode(&nodes, el.N2, el.N3)
			N7  = midNode(&nodes, el.N3, el.N4)
			N8  = midNode(&nodes, el.N1, el.N4)
			N9  = midNode(&nodes, el.N1, el.N3)
			N10 = midNode(&nodes, el.N2, el.N4)
			N11 = midNode(&nodes, el.N1, el.N2, el.N3, el.N4)
		)

		newElements = append(newElements,
			Element{el.N1, N5, N8, N9},
			Element{el.N4, N7, N8, N10},
			Element{el.N3, N6, N7, N9},
			Element{el.N2, N5, N6, N10},
			Element{N11, N5, N8, N9},
			Element{N11, N7, N8, N10},
			Element{N11, N6, N7, N9},
			Element{N11, N5, N6, N10},
			Element{N11, N5, N8, N10},
			Element{N11, N7, N8, N9},
			Element{N11, N6, N7, N10},
			Element{N11, N5, N6, N9},
		)

		if extraCount == 3 {
			boundaryNodeIndexes = append(boundaryNodeIndexes, N5, N6, N9)
		}
	}

	return newElements, nodes, boundaryNodeIndexes
}

func inArr(arr []int, el int) bool {
	if ind := sort.SearchInts(arr, el); ind == len(arr) || arr[ind] != el {
		return false
	}

	return true
}

func show3D(nodes []Node, elements []Element) {
	// Creates window and OpenGL context
	wmgr, err := window.Manager("glfw")
	if err != nil {
		panic(err)
	}
	win, err := wmgr.CreateWindow(800, 600, "Hello G3N", false)
	if err != nil {
		panic(err)
	}

	// OpenGL functions must be executed in the same thread where
	// the context was created (by CreateWindow())
	runtime.LockOSThread()

	// Create OpenGL state
	gs, err := gls.New()
	if err != nil {
		panic(err)
	}

	// Sets the OpenGL viewport size the same as the window size
	// This normally should be updated if the window is resized.
	width, height := win.Size()
	gs.Viewport(0, 0, int32(width), int32(height))

	// Creates scene for 3D objects
	scene := core.NewNode()

	// Adds white ambient light to the scene
	ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	scene.Add(ambLight)

	// Adds a perspective camera to the scene
	// The camera aspect ratio should be updated if the window is resized.
	aspect := float32(width) / float32(height)

	camera := camera.NewPerspective(35, aspect, 0.01, 1000)
	camera.SetPosition(6, 5, 6)
	//camera.SetDirection(0-4, 0-4, -4)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(2)
	scene.Add(axis)

	// Creates a wireframe sphere positioned at the center of the scene
	//geom := geometry.NewSphere(2, 16, 16, 0, math.Pi*2, 0, math.Pi)
	//geom := geometry.NewBox(4, 4, 4)
	geom := NewGeometry(nodes, elements)
	mat := material.NewStandard(math32.NewColor("White"))
	mat.SetSide(material.SideDouble)
	mat.SetWireframe(true)
	sphere := graphic.NewMesh(geom, mat)

	scene.Add(sphere)

	// Creates a renderer and adds default shaders
	rend := renderer.NewRenderer(gs)
	err = rend.AddDefaultShaders()
	if err != nil {
		panic(err)
	}
	rend.SetScene(scene)

	// Sets window background color
	gs.ClearColor(0, 0, 0, 1.0)

	// Render loop
	for !win.ShouldClose() {

		// Rotates the sphere a bit around the Z axis (up)
		sphere.AddRotationY(0.001)

		// Render the scene using the specified camera
		rend.Render(camera)

		// Update window and checks for I/O events
		win.SwapBuffers()
		wmgr.PollEvents()
	}
}

func midNode(nodes *[]Node, midNodes ...int) int {
	var node Node
	for _, n := range midNodes {
		node.X += (*nodes)[n].X
		node.Y += (*nodes)[n].Y
		node.Z += (*nodes)[n].Z
	}
	node.X /= float64(len(midNodes))
	node.Y /= float64(len(midNodes))
	node.Z /= float64(len(midNodes))

	*nodes = append(*nodes, node)

	return len(*nodes) - 1
}
