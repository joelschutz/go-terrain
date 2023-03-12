package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/joelschutz/go-terrain/pkg"
	"github.com/joelschutz/go-terrain/util"
)

const (
	meshSize   = 512
	mapScale   = 40
	iterations = 200
)

var (
	age       = 0
	rainCount = meshSize * 10
)

func GeneratePoints(heightMap pkg.HeightMap) (pkg.RainSim, *core.Node) {
	sim := pkg.RainSim{}
	simPoints := core.NewNode()
	for i := 0; i < rainCount; i++ {
		p := pkg.NewParticle(mgl32.Vec2{rand.Float32() * meshSize, rand.Float32() * meshSize}, mgl32.Vec2{}, 1, 1, 1)
		sim = append(sim, p)
		pg := graphic.NewMesh(geometry.NewSphere(1, 4, 4), material.NewStandard(math32.NewColor("blue")))
		vec := math32.NewVec3()
		vec.Set(p.X()-meshSize/2, p.Y()-meshSize/2, float32(heightMap[int(p.X())][int(p.Y())]*mapScale))
		pg.SetPositionVec(vec)
		simPoints.Add(pg)
	}
	return sim, simPoints
}

func G3nApp() {
	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 200)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// // Create Noise
	heightMap := pkg.NewPerlinMap(640, 100)
	normalMap := pkg.NewNormalMap(heightMap, 100)

	sim, simPoints := GeneratePoints(heightMap)
	scene.Add(simPoints)

	// Create a blue torus and add it to the scene
	geom := geometry.NewSegmentedPlane(meshSize, meshSize, meshSize, meshSize)
	geom.OperateOnVertices(func(vertex *math32.Vector3) bool {
		vertex.SetZ(float32(heightMap[int(vertex.X+meshSize/2)][int(vertex.Y+meshSize/2)] * mapScale))
		return false
	})
	i := 0
	nmTruc := normalMap.Truncate()
	geom.OperateOnVertexNormals(func(normal *math32.Vector3) bool {
		normal.X = float32(nmTruc[i].X())
		normal.Y = float32(nmTruc[i].Y())
		normal.Z = float32(nmTruc[i].Z())
		i++
		return false
	})
	tex1 := texture.NewTexture2DFromRGBA(util.AsRGBA(heightMap))
	mat := material.NewStandard(math32.NewColor("white"))
	mat.AddTexture(tex1)
	mat.SetDepthMask(true)
	// mat.SetWireframe(true)
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Create and add a button to the scene
	btn := gui.NewButton("Update")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		go func() {
			for i := 0; i < iterations; i++ {
				sim.Update(&heightMap, true)
			}
		}()
	})
	scene.Add(btn)

	btn2 := gui.NewButton("Reset")
	btn2.SetPosition(100, 80)
	btn2.SetSize(40, 40)
	btn2.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		mesh.GetGeometry().OperateOnVertices(func(vertex *math32.Vector3) bool {
			vertex.SetZ(float32(heightMap[int(vertex.X+meshSize/2)][int(vertex.Y+meshSize/2)] * mapScale))
			return false
		})

		scene.Remove(simPoints)
		sim, simPoints = GeneratePoints(heightMap)
		scene.Add(simPoints)
	})
	scene.Add(btn2)

	scene.Add(light.NewAmbient(&math32.Color{1, 1, 1}, 1))

	// Helpers
	scene.Add(helper.NewAxes(0.5))
	lb := gui.NewLabel(fmt.Sprintf("X:%v Y:%v Z%v", cam.Position().X, cam.Position().Y, cam.Position().Z))
	scene.Add(lb)
	// scene.Add(helper.NewNormals(mesh, 10, &math32.Color{0, 1, 0}, .5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
		lb.SetText(fmt.Sprintf("X:%v Y:%v Z%v", cam.Position().X, cam.Position().Y, cam.Position().Z))
		for i, n := range simPoints.Children() {
			pg := n.(*graphic.Mesh)
			p := sim[i]
			pg.SetPosition(p.X()-meshSize/2, p.Y()-meshSize/2, float32(heightMap[int(p.X())][int(p.Y())]*mapScale))
			if p.Stopped {
				pg.SetMaterial(material.NewStandard(math32.NewColor("red")))
			}
		}
		age++
	})
}
