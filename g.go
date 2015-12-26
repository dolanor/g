package main

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
	"azul3d.org/mouse.v1"
	"fmt"
	"image"
	"log"
)

var glslVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec2 TexCoord0;

uniform mat4 MVP;

varying vec2 tc0;

void main()
{
	tc0 = TexCoord0;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
`)

var glslFrag = []byte(`
#version 120

varying vec2 tc0;

uniform sampler2D Texture0;
uniform bool BinaryAlpha;

void main()
{
	gl_FragColor = texture2D(Texture0, tc0);
	if(BinaryAlpha && gl_FragColor.a < 0.5) {
		discard;
	}
}
`)

func gfxLoop(w window.Window, r gfx.Renderer) {

	cam := gfx.NewCamera()
	camFOV := 75.0
	camNear := 0.0001
	camFar := 1000.0

	cam.SetPersp(r.Bounds(), camFOV, camNear, camFar)
	cam.SetPos(lmath.Vec3{0, -2, 0})

	rtColor := gfx.NewTexture()
	rtColor.MinFilter = gfx.LinearMipmapLinear
	rtColor.MagFilter = gfx.Linear

	cfg := r.GPUInfo().RTTFormats.ChooseConfig(gfx.Precision{
		RedBits: 8, BlueBits: 8, GreenBits: 8}, true)

	cfg.Color = rtColor
	cfg.Bounds = image.Rect(0, 0, 512, 512)

	rtCanvas := r.RenderToTexture(cfg)

	if rtCanvas == nil {
		log.Fatal("Graphic card doesn't support render to texture")
	}

	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

	cardMesh := gfx.NewMesh()
	cardMesh.Vertices = []gfx.Vec3{
		{-1, 0, -1},
		{1, 0, -1},
		{-1, 0, 1},

		{-1, 0, 1},
		{1, 0, -1},
		{1, 0, 1},
	}

	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: []gfx.TexCoord{
				{0, 1},
				{1, 1},
				{0, 0},

				{0, 0},
				{1, 1},
				{1, 0},
			},
		},
	}

	card := gfx.NewObject()
	card.FaceCulling = gfx.NoFaceCulling
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{rtColor}
	card.Meshes = []*gfx.Mesh{cardMesh}

	go func() {
		events := make(chan window.Event, 256)
		w.Notify(events, window.AllEvents)

		for event := range events {
			switch e := event.(type) {
			case keyboard.TypedEvent:
				fmt.Println("Pressed", event)
			case keyboard.StateEvent:
				trans := lmath.Vec3{}
				switch e.Key {
				case keyboard.ArrowLeft:
					trans.X -= .1
				case keyboard.ArrowRight:
					trans.X += .1
				case keyboard.ArrowUp:
					trans.Z += .1
				case keyboard.ArrowDown:
					trans.Z -= .1
				}
				card.SetPos(lmath.Vec3{
					card.Pos().X + trans.X,
					0,
					card.Pos().Z + trans.Z,
				})

			case window.FramebufferResized:
				cam.Lock()
				cam.SetPersp(r.Bounds(), camFOV, camNear, camFar)
				cam.Unlock()
			}
		}
	}()

	stripeColor1 := gfx.Color{1, 0, 0, 1}
	stripeColor2 := gfx.Color{0, 0, 1, 1}

	stripeWidth := 12
	flipColor := false

	b := rtCanvas.Bounds()

	for i := 0; (i * stripeWidth) < b.Dx(); i++ {
		flipColor = !flipColor
		x := i * stripeWidth
		dst := image.Rect(x, b.Min.Y, x+stripeWidth, b.Max.Y)
		if flipColor {
			rtCanvas.Clear(dst, stripeColor1)
		} else {
			rtCanvas.Clear(dst, stripeColor2)
		}
	}
	rtCanvas.Render()

	for {
		rot := card.Rot()
		card.SetRot(lmath.Vec3{
			X: rot.X,
			Y: rot.Y,
			Z: rot.Z + (15 * r.Clock().Dt()),
		})
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// The keyboard is monitored for you, simply check if a key is down:
		if w.Keyboard().Down(keyboard.Space) {
			// Clear a red rectangle.
			r.Clear(image.Rect(0, 0, 100, 100), gfx.Color{1, 0, 0, 1})
		}

		// And the same thing with the mouse, check if a mouse button is down:
		if w.Mouse().Down(mouse.Left) {
			// Clear a blue rectangle.
			r.Clear(image.Rect(100, 100, 200, 200), gfx.Color{0, 0, 1, 1})
		}

		r.Draw(image.Rect(0, 0, 0, 0), card, cam)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
