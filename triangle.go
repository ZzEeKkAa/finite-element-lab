// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

func NewGeometry(nodes []Node, elements []Element) *geometry.Geometry {
	geometry := new(geometry.Geometry)
	geometry.Init()

	positions := math32.NewArrayF32(0, len(nodes))
	indices := math32.NewArrayU32(0, len(elements)*4)

	k := 4.
	for _, n := range nodes {
		positions.Append(float32(k*(n.X-0.5)), float32(k*n.Z), float32(k*(n.Y-0.5)))
	}

	for _, e := range elements {
		indices.Append(uint32(e.N1), uint32(e.N2), uint32(e.N3))
		indices.Append(uint32(e.N1), uint32(e.N2), uint32(e.N4))
		indices.Append(uint32(e.N1), uint32(e.N3), uint32(e.N4))
		indices.Append(uint32(e.N2), uint32(e.N3), uint32(e.N4))
	}

	geometry.SetIndices(indices)
	geometry.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))

	return geometry
}
