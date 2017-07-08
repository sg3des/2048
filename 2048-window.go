package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/sg3des/fizzgui"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"
)

var (
	window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider

	fontfilename = "Roboto-Black.ttf"

	TextFont      *fizzgui.Font
	TextFontSmall *fizzgui.Font
	NumsFont      *fizzgui.Font
)

func NewWindow(title string, w, h int) error {
	runtime.LockOSThread()

	window, gfx = initGraphics(title, w, h)

	err := fizzgui.Init(window, gfx)
	if err != nil {
		return fmt.Errorf("Failed initialize fizzgui, reason: %s", err)
	}

	TextFont, err = fizzgui.NewFont("Default", fontfilename, 28, fizzgui.FontGlyphs)
	if err != nil {
		return fmt.Errorf("Failed to load the Text font, reason: %s", err)
	}

	TextFontSmall, err = fizzgui.NewFont("Text", fontfilename, 20, fizzgui.FontGlyphs)
	if err != nil {
		return fmt.Errorf("Failed to load the Text font, reason: %s", err)
	}

	//load a default font
	NumsFont, err = fizzgui.NewFont("Nums", fontfilename, 41, "012345689")
	if err != nil {
		return fmt.Errorf("Failed to load the Default font, reason: %s", err)
	}

	return nil
}

func RenderLoop() {
	for {
		t := time.Now()

		w, h := window.GetFramebufferSize()
		gfx.Viewport(0, 0, int32(w), int32(h))
		gfx.ClearColor(0.9, 0.9, 0.9, 1)
		gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		// draw the user interface
		fizzgui.Construct()

		// draw the screen and get any input
		window.SwapBuffers()
		glfw.PollEvents()

		dt := float32(time.Now().Sub(t).Seconds())
		Transitions(dt)

		if window.ShouldClose() {
			Close()
		}
	}
}

// initGraphics creates an OpenGL window and initializes the required graphics libraries.
// It will either succeed or panic.
func initGraphics(title string, w int, h int) (*glfw.Window, graphicsprovider.GraphicsProvider) {

	err := glfw.Init()
	if err != nil {
		panic("Can't init glfw! " + err.Error())
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 0)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	window, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		panic("Failed to create the main window! " + err.Error())
	}

	window.MakeContextCurrent()

	glfw.SwapInterval(1) // if 0 disable v-sync

	// initialize OpenGL
	gfx, err := opengl.InitOpenGL()
	if err != nil {
		panic("Failed to initialize OpenGL! " + err.Error())
	}
	fizzle.SetGraphics(gfx)

	// set some additional OpenGL flags
	gfx.BlendEquation(graphicsprovider.FUNC_ADD)
	gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)
	gfx.Enable(graphicsprovider.BLEND)
	gfx.Enable(graphicsprovider.TEXTURE_2D)
	gfx.Enable(graphicsprovider.CULL_FACE)

	window.SetKeyCallback(keyCallback)

	return window, gfx
}
