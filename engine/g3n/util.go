package engine

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

//
// PointLightMesh
//
type PointLightMesh struct {
	*graphic.Mesh
	Light *light.Point
}

func (pl *PointLightMesh) SetPosition(x, y, z float32) {
	pl.Light.SetPosition(x, y, z)
	pl.Mesh.SetPosition(x, y, z)

}

func NewPointLightMesh(color *math32.Color, intensity float32) *PointLightMesh {

	l := new(PointLightMesh)

	geom := geometry.NewSphere(1, 5, 5)
	mat := material.NewStandard(color)
	mat.SetUseLights(0)
	mat.SetEmissiveColor(color)
	l.Mesh = graphic.NewMesh(geom, mat)
	l.Mesh.SetVisible(true)

	l.Light = light.NewPoint(color, intensity)
	l.Light.SetLinearDecay(1)
	l.Light.SetQuadraticDecay(1)
	l.Light.SetVisible(true)

	l.SetPosition(0, 0, 10)
	l.Mesh.Add(l.Light)

	return l
}

//
// SpotLightMesh
//
type SpotLightMesh struct {
	*graphic.Mesh
	Light *light.Spot
}

func NewSpotLightMesh(color *math32.Color) *SpotLightMesh {

	l := new(SpotLightMesh)

	geom := geometry.NewCylinder(32, 16, 16, 16, true, true)
	mat1 := material.NewStandard(color)
	mat2 := material.NewStandard(color)
	mat2.SetEmissiveColor(color)
	l.Mesh = graphic.NewMesh(geom, nil)
	l.Mesh.AddGroupMaterial(mat1, 0)
	l.Mesh.AddGroupMaterial(mat2, 1)

	l.Light = light.NewSpot(color, 2.0)
	l.Light.SetDirection(0, 0, -1)
	l.Light.SetCutoffAngle(45)
	l.Light.SetLinearDecay(0)
	l.Light.SetQuadraticDecay(0)

	l.Mesh.Add(l.Light)

	return l
}
