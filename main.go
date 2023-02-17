package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"net"
	"os"
)

const (
	serverProtocol = "tcp"
	serverHost     = ""
	serverPort     = "8080"
)

func ClientSend(i int) {
	connection, err := net.Dial(serverProtocol, serverHost+":"+serverPort)

	if err != nil {
		panic(err)
	}

	_, err = connection.Write([]byte("Hello Server! Greetings."))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("[CLIENT] Error reading:", err.Error())
	}

	fmt.Println(i, "Received: ", string(buffer[:mLen]))

	defer connection.Close()
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("[SERVER] Error reading:", err.Error())
	}

	fmt.Println("Received: ", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
}

// socket-server project main.go
func ServerLoop() {
	fmt.Println("Server Running...")

	server, err := net.Listen(serverProtocol, serverHost+":"+serverPort)

	if err != nil {
		fmt.Println("Error listening:", err.Error())

		os.Exit(1)
	}

	defer server.Close()

	fmt.Println("Listening on " + serverHost + ":" + serverPort)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("client connected")

		go processClient(connection)
	}
}

func AllTogether() {
	go ServerLoop()

	var i int = 1

	timer := Timer{
		config:  SIMPLE | LOOP,
		delay:   5000,
		current: sdl.GetTicks(),
		state:   NONE,
	}

	for {
		timer.Update()

		if timer.state&DONE > 0 {
			go ClientSend(i)
			i++
		}
	}

}

type Texture struct {
	id      uint32
	w, h    int32
	surface *sdl.Surface
}

var defaultTexture Texture

func loadXPM(filepath string) Texture {
	var id uint32 = 0

	// file, err := os.ReadFile(filepath)

	// if err != nil {
	// panic(err)
	// }

	// fmt.Println(string(file))

	regsurface, err := img.Load(filepath)

	if err != nil {
		panic(err)
	}

	surface, err := regsurface.ConvertFormat(uint32(sdl.PIXELFORMAT_RGBA32), 0)

	if err != nil {
		panic(err)
	}

	var mode uint32 = gl.RGB
	var internal_format int32 = gl.RGBA

	if surface.Format.BytesPerPixel == 4 {
		mode = gl.RGBA
	}

	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)

	gl.TexImage2D(gl.TEXTURE_2D, 0, internal_format, surface.W, surface.H, 0,
		mode, gl.UNSIGNED_BYTE, gl.Ptr(surface.Pixels()))

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	regsurface.Free()
	// surface.Free()

	return Texture{w: surface.W, h: surface.H, id: id, surface: surface}
}

type Ticks struct {
	frames, fps                                  uint32
	previous, current, delta, frame, accumulator float32
}

var ticks Ticks

const NONE = 0

const (
	START = 1
	WAIT  = 2
	DONE  = 4
)

const (
	SIMPLE  = 1
	TWO_WAY = 2
	LOOP    = 4
)

type Timer struct {
	state, config, delay, restartDelay, current uint32
}

func (t *Timer) Set(state uint32) {
	t.state = state
	t.current = sdl.GetTicks()
}

func (t *Timer) Update() {
	diff := sdl.GetTicks() - t.current

	if t.state == NONE {
		t.Set(START)
	} else if t.state&DONE > 0 {
		if t.config&LOOP > 0 {
			t.Set(START)
		}
	} else if diff >= t.delay {
		t.Set(DONE)
	}
}

func InitTicks() {
	ticks = Ticks{
		frames: 0, fps: 0,
		previous: 0, delta: 0.01, frame: 0, accumulator: 0,
		current: float32(sdl.GetTicks() / 1000),
	}
}

func (t Ticks) dt() float32 {
	return t.delta
}

func (t *Ticks) Update() {
	t.frames++

	t.previous = t.current
	t.current = float32(sdl.GetTicks() / 1000)
	t.frame = t.current - t.previous

	if t.frame > 0.25 {
		t.frame = 0.25
	}

	t.accumulator += t.frame

	t.fps = t.frames / (sdl.GetTicks() / 1000)
}

type vec2 struct {
	x, y float32
}

func (a vec2) Add(b vec2) vec2 {
	return vec2{a.x + b.x, a.y + b.y}
}

func (a vec2) Sub(b vec2) vec2 {
	return vec2{a.x - b.x, a.y - b.y}
}

func (a vec2) Mul(b vec2) vec2 {
	return vec2{a.x * b.x, a.y * b.y}
}

func (a *vec2) AddF(b float32) {
	a.x += b
	a.y += b
}

func (a vec2) Div(b vec2) vec2 {
	return vec2{a.x / b.x, a.y / b.y}
}

type vec3 struct {
	x, y, z float32
}

func (a vec3) Add(b vec3) vec3 {
	return vec3{a.x + b.x, a.y + b.y, a.z + b.z}
}

func (a vec3) Sub(b vec3) vec3 {
	return vec3{a.x - b.x, a.y - b.y, a.z - b.z}
}

func (a vec3) Mul(b vec3) vec3 {
	return vec3{a.x * b.x, a.y * b.y, a.z * b.z}
}

func (a vec3) Div(b vec3) vec3 {
	return vec3{a.x / b.x, a.y / b.y, a.z / b.z}
}

type vec4 struct {
	x, y, z, w float32
}

func (a vec4) Add(b vec4) vec4 {
	return vec4{a.x + b.x, a.y + b.y, a.z + b.z, a.w + b.w}
}

func (a vec4) Sub(b vec4) vec4 {
	return vec4{a.x - b.x, a.y - b.y, a.z - b.z, a.w - b.w}
}

func (a vec4) Mul(b vec4) vec4 {
	return vec4{a.x * b.x, a.y * b.y, a.z * b.z, a.w * b.w}
}

func (a vec4) Div(b vec4) vec4 {
	return vec4{a.x / b.x, a.y / b.y, a.z / b.z, a.w / b.w}
}

type Mat4 struct {
	data [4][4]float32
}

func Identity() Mat4 {
	return Mat4{
		[4][4]float32{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func (m1 Mat4) Mul(m2 Mat4) Mat4 {
	curr, result := float32(0), Mat4{}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			curr = 0

			for k := 0; k < 4; k++ {
				curr += (m1.data[i][k] * m2.data[k][j])
			}

			result.data[i][j] = curr
		}
	}

	return result
}

func (m1 Mat4) ToArray() []float32 {
	result := []float32{}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result = append(result, m1.data[i][j])
		}
	}

	return result
}

func Ortho(W float32, H float32) Mat4 {
	r, t := W, float32(0)
	l, b := float32(0), H
	f, n := float32(1), float32(-1)

	matrix := Identity()

	matrix.data[0][0] = 2 / (r - l)
	matrix.data[0][3] = -((r + l) / (r - l))

	matrix.data[1][1] = 2 / (t - b)
	matrix.data[1][3] = -((t + b) / (t - b))

	matrix.data[2][2] = -2 / (f - n)
	matrix.data[2][3] = -((f + n) / (f - n))

	return matrix
}

func (m1 Mat4) Translate(pos vec2) Mat4 {
	transMatrix := Identity()

	// {1, 0, 0, pos.x}
	// {0, 1, 0, pos.y}
	// {0, 0, 1, 0    }
	// {0, 0, 0, 1    }

	transMatrix.data[0][3] = pos.x
	transMatrix.data[1][3] = pos.y

	return m1.Mul(transMatrix)
}

func (m Mat4) Scale(size vec3) Mat4 {
	scaleMatrix := Identity()

	// {size.x, 0,      0,      0}
	// {0,      size.y, 0,      0}
	// {0,      0,      size.z, 0}
	// {0,      0,      0,      1}

	scaleMatrix.data[0][0] = size.x
	scaleMatrix.data[1][1] = size.y
	scaleMatrix.data[2][2] = size.z

	return m.Mul(scaleMatrix)
}

func (s Shader) Location(str string) int32 {
	return gl.GetUniformLocation(s.id, gl.Str(str+"\x00"))
}

func (s Shader) SetInt(str string, v int32) {
	gl.Uniform1i(s.Location(str), v)
}

func (s Shader) SetVec4(str string, v vec4) {
	raw := [4]float32{v.x, v.y, v.z, v.w}

	gl.Uniform4fv(s.Location(str), 1, (*float32)(gl.Ptr(&raw[0])))
}

func (s Shader) SetMat4(str string, m Mat4) {
	gl.UniformMatrix4fv(s.Location(str), 1, true, (*float32)(gl.Ptr(m.ToArray())))
}

func (s Shader) Use() {
	gl.Flush()

	gl.UseProgram(s.id)
	gl.BindVertexArray(s.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
}

func getModel(pos vec2, size vec2) Mat4 {
	m := Identity()

	m = m.Translate(pos)
	m = m.Scale(vec3{size.x, size.y, 0})

	return m
}

const W, H = 320, 180

type Shader struct {
	id, vId, fId  uint32
	vao, vbo, ebo uint32
}

func LoadShaderFile(filepath string, shaderType uint32) uint32 {
	file, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	id := gl.CreateShader(shaderType)

	source, free := gl.Strs(string(file) + "\x00")

	defer free()

	gl.ShaderSource(id, 1, source, nil)
	gl.CompileShader(id)

	isCompiled := int32(0)

	gl.GetShaderiv(id, gl.COMPILE_STATUS, &isCompiled)

	if isCompiled == gl.FALSE {

		maxLength := int32(0)
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &maxLength)

		// The maxLength includes the NULL character
		var errorLog []byte = make([]byte, int(maxLength))
		gl.GetShaderInfoLog(id, maxLength, &maxLength, &errorLog[0])

		gl.DeleteShader(id) // Don't leak the shader.

		panic(string(errorLog))
	}

	return id
}

func LoadShader(vertexFilePath string, fragmentFilePath string) Shader {
	s := Shader{
		gl.CreateProgram(),
		LoadShaderFile(vertexFilePath, gl.VERTEX_SHADER),
		LoadShaderFile(fragmentFilePath, gl.FRAGMENT_SHADER),
		0, 0, 0,
	}

	gl.AttachShader(s.id, s.vId)
	gl.AttachShader(s.id, s.fId)

	gl.LinkProgram(s.id)

	isLinked := int32(0)

	gl.GetProgramiv(s.id, gl.LINK_STATUS, &isLinked)

	if isLinked == gl.FALSE {

		maxLength := int32(0)
		gl.GetShaderiv(s.id, gl.INFO_LOG_LENGTH, &maxLength)

		var errorLog []byte = make([]byte, int(maxLength))
		gl.GetShaderInfoLog(s.id, maxLength, &maxLength, &errorLog[0])

		gl.DeleteShader(s.id) // Don't leak the shader.

		panic(string(errorLog))
	}

	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &s.vbo)
	gl.GenBuffers(1, &s.ebo)

	return s
}

const (
	bit8  = 1
	bit16 = 2
	bit32 = 4
	bit64 = 8
)

func DefaultShader() Shader {
	s := LoadShader("vertex.glsl", "fragment.glsl")

	s.Use()

	points := []float32{
		0, 0, //
		0, 1, //
		1, 1, //
		1, 0, //
	}

	indices := []int32{
		0, 1, 2, //
		0, 2, 3, //
	}

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*bit32, nil)
	gl.EnableVertexAttribArray(0)

	gl.BufferData(gl.ARRAY_BUFFER, len(points)*bit32, gl.Ptr(points), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*bit32, gl.Ptr(indices), gl.STATIC_DRAW)

	s.SetInt("uType", 0)
	s.SetInt("uImage", 0)

	s.SetVec4("uColor", vec4{1, 1, 1, 1})

	s.SetMat4("uModel", Identity())
	s.SetMat4("uProjection", Ortho(W, H))

	return s
}

var engine Engine
var defaultShader Shader

type ButtonCallback func(btn *Button, data any)

type ButtonData struct {
	pos, size vec2
}

type Button interface {
	Draw()
	Update()
	IsClicked(spot vec2) bool
	Hover()
	UnHover()
}

type SimpleButton struct {
	ButtonData
	toggle, border bool
	color, sprite  vec4
}

func (b SimpleButton) IsClicked(spot vec2) bool {
	if spot.x >= b.pos.x && spot.x <= b.pos.x+b.size.x &&
		spot.y >= b.pos.y && spot.y <= b.pos.y+b.size.y {
		return true
	} else {
		return false
	}
}

func (b *SimpleButton) Hover() {
	b.border = true
}

func (b *SimpleButton) UnHover() {
	b.border = false
}

type Player struct {
	id            uint8
	pos           vec2
	color, sprite vec4
}

var player1, player2 Player
var activePlayer *Player = nil

func CheckWinCondition() {

}

func (b *SimpleButton) Update() {
	if !b.toggle {
		b.toggle = true
		b.color = vec4{0, 1, 0, 1}
		b.sprite = activePlayer.sprite

		if activePlayer == &player1 {
			activePlayer = &player2
		} else {
			activePlayer = &player1
		}
	}
}

func (t Texture) Coords(coords vec4) vec4 {
	return coords.Div(vec4{float32(t.w), float32(t.h), float32(t.w),
		float32(t.h)})
}

func (b SimpleButton) Draw() {
	if b.border {
		defaultShader.SetMat4("uModel", getModel(b.pos.Sub(vec2{2, 2}),
			b.size.Add(vec2{4, 4})))

		defaultShader.SetVec4("uOffset", defaultTexture.Coords(vec4{0, 16, 16, 16}))
		defaultShader.SetVec4("uColor", vec4{.5, .5, 0, 1})

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	}

	defaultShader.SetMat4("uModel", getModel(b.pos, b.size))

	defaultShader.SetVec4("uOffset", defaultTexture.Coords(vec4{0, 16, 16, 16}))
	defaultShader.SetVec4("uColor", vec4{.2, .2, .2, 1})

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	defaultShader.SetVec4("uOffset", b.sprite)
	defaultShader.SetVec4("uColor", b.color)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

type Engine struct {
	run     bool
	window  *sdl.Window
	context sdl.GLContext

	buttons []Button

	w, h         int32
	realW, realH int32
}

func (engine *Engine) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	var flags uint32 = sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL
	// | sdl.WINDOW_FULLSCREEN_DESKTOP

	window, err := sdl.CreateWindow("gocard", 0, 0, W, H, flags)

	if err != nil {
		panic(err)
	}

	context, err := window.GLCreateContext()

	if err != nil {
		panic(err)
	}

	if err = gl.Init(); err != nil {
		panic(err)
	}

	// sdl.GLSetSwapInterval(0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	*engine = Engine{true, window, context, []Button{}, W, H, W, H}
}

func CheckButtonPress(spot vec2, buttons []Button) {
	factor := vec2{
		float32(engine.realW / engine.w),
		float32(engine.realH / engine.h),
	}

	spot = spot.Div(factor)

	for i := range buttons {
		if buttons[i].IsClicked(spot) {
			buttons[i].Update()
		}
	}
}

func (engine *Engine) Physics() {
	for ticks.accumulator >= ticks.dt() {
		// we gotta do something
		ticks.accumulator -= ticks.dt()
	}
}

func (engine *Engine) Event() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.WindowEvent:
			switch t.Event {
			case sdl.WINDOWEVENT_RESIZED:
				engine.realW, engine.realH = t.Data1, t.Data2
				gl.Viewport(0, 0, t.Data1, t.Data2)
			}
			break
		case *sdl.QuitEvent:
			engine.run = false
			break
		case *sdl.MouseButtonEvent:
			switch t.Type {
			case sdl.MOUSEBUTTONDOWN:
				switch t.Button {
				case sdl.BUTTON_LEFT:
					CheckButtonPress(vec2{float32(t.X), float32(t.Y)},
						engine.buttons)
					break
				}
				break
			}
			break
		case *sdl.MouseMotionEvent:
			activePlayer.pos = vec2{
				float32(t.X),
				float32(t.Y),
			}

			// activePlayer.pos = activePlayer.pos.Add(vec2{
			// 	float32(t.XRel),
			// 	float32(t.YRel),
			// })

			spot := activePlayer.pos

			factor := vec2{
				float32(engine.realW / engine.w),
				float32(engine.realH / engine.h),
			}

			spot = spot.Div(factor)

			for i := range engine.buttons {
				if engine.buttons[i].IsClicked(spot) {
					engine.buttons[i].Hover()
				} else {
					engine.buttons[i].UnHover()
				}
			}

			break
		case *sdl.KeyboardEvent:
			if t.Repeat == 0 && t.Type == sdl.KEYDOWN {
				switch t.Keysym.Sym {
				case sdl.K_a:
					if engine.window.GetFlags()&sdl.WINDOW_FULLSCREEN_DESKTOP > 0 {
						engine.window.SetFullscreen(0)
					} else {
						engine.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					}

					break
				case sdl.K_s:
					if sdl.GetRelativeMouseMode() {
						sdl.SetRelativeMouseMode(false)
					} else {
						sdl.SetRelativeMouseMode(true)
					}
					break
				}
			}

			break
		}
	}
}

func AddGridButtons(k float32, l float32, squareSize float32) []Button {
	var buttons []Button

	size := vec2{squareSize, squareSize}

	middle := vec2{
		float32(engine.w)/2 - ((k * (size.x * 1.5)) / 2),
		float32(engine.h)/2 - ((l * (size.y * 1.5)) / 2),
	}

	for i := float32(0); i < k; i++ {
		for j := float32(0); j < l; j++ {
			buttons = append(buttons, &SimpleButton{
				ButtonData: ButtonData{
					pos:  middle.Add(vec2{i * (size.x * 1.5), j * (size.y * 1.5)}),
					size: size,
				},
				border: false,
				toggle: false,
				color:  vec4{0, 0, 1, 1},
			})
		}
	}

	return buttons
}

func Norm() vec2 {
	return vec2{
		float32(engine.realW / engine.w),
		float32(engine.realH / engine.h),
	}
}

func (p Player) Draw() {
	pos := p.pos.Sub(vec2{4, 4}.Mul(Norm()))

	defaultShader.SetMat4("uModel", getModel(pos.Div(Norm()),
		vec2{8, 8}))

	defaultShader.SetVec4("uOffset", defaultTexture.Coords(vec4{32, 0, 16, 16}))
	defaultShader.SetVec4("uColor", p.color)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func main() {
	go AllTogether()

	InitTicks()
	engine.Init()

	defer sdl.Quit()

	defaultShader = DefaultShader()

	defaultShader.SetInt("uType", 1)
	defaultShader.SetMat4("uModel", getModel(vec2{32, 32}, vec2{32, 32}))
	defaultShader.SetVec4("uOffset", vec4{0, 0, .1, .1})

	defaultTexture = loadXPM("spritesheet.png")

	engine.buttons = AddGridButtons(3, 3, 16)

	timer := Timer{
		state:   NONE,
		config:  SIMPLE | LOOP,
		delay:   5000,
		current: sdl.GetTicks(),
	}

	player1 = Player{
		id:     0,
		pos:    vec2{W / 2, H / 2},
		sprite: defaultTexture.Coords(vec4{0, 0, 16, 16}),
		color:  vec4{0, 1, 0, 1},
	}

	player2 = Player{
		id:     1,
		pos:    vec2{W / 2, H / 2},
		sprite: defaultTexture.Coords(vec4{16, 0, 16, 16}),
		color:  vec4{1, 0, 0, 1},
	}

	activePlayer = &player1

	for engine.run {
		ticks.Update()
		engine.Event()
		engine.Physics()

		timer.Update()

		if timer.state&DONE > 0 {
			fmt.Println("fps: ", ticks.fps)
		}

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		for i := range engine.buttons {
			engine.buttons[i].Draw()
		}

		activePlayer.Draw()

		engine.window.GLSwap()
	}

	fmt.Println("goodbye cardgame")
}
