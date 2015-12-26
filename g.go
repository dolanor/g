package main

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
	"fmt"
	"image"
)

func gfxLoop(w window.Window, r gfx.Renderer) {

	go func() {
		events := make(chan window.Event, 256)
		w.Notify(events, window.AllEvents)

		for event := range events {
			fmt.Println("Pressed", event)
		}
	}()

	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

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

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
