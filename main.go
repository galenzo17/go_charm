package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// Estilos para los colores neón
var (
	neonPink   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF10F0"))
	neonBlue   = lipgloss.NewStyle().Foreground(lipgloss.Color("#10F0FF"))
	neonGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10FF50"))
	neonYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF10"))
	neonColors = []*lipgloss.Style{&neonPink, &neonBlue, &neonGreen, &neonYellow}
)

// Particle representa una partícula que sigue al cursor
type Particle struct {
	x, y           float64
	targetX, targetY float64
	xVel, yVel     float64
	springX, springY harmonica.Spring
	char            string
	style           *lipgloss.Style
}

type model struct {
	width, height   int
	cursorX, cursorY int
	particles       []Particle
	frameCount      int
	trail           [][2]int
	maxTrail        int
	clickX, clickY  int
	clickActive     bool
	clickRadius     int
	clickMaxRadius  int
	clickStyle      *lipgloss.Style
}

func initialModel() model {
	m := model{
		width:         80,
		height:        24,
		maxTrail:      20,
		clickMaxRadius: 10,
	}

	// Partículas que siguen al cursor
	m.particles = make([]Particle, 8)
	for i := range m.particles {
		frequency := 4.0 + float64(i)*0.5
		damping := 0.3 + float64(i)*0.05

		m.particles[i] = Particle{
			x:       float64(m.width / 2),
			y:       float64(m.height / 2),
			targetX: float64(m.width / 2),
			targetY: float64(m.height / 2),
			springX: harmonica.NewSpring(harmonica.FPS(30), frequency, damping),
			springY: harmonica.NewSpring(harmonica.FPS(30), frequency, damping),
			char:    "●",
			style:   neonColors[i%len(neonColors)],
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.ClearScreen,
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}
type mouseMsg struct{ x, y int }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.MouseMsg:
		m.cursorX, m.cursorY = msg.X, msg.Y

		// Registra el rastro del cursor
		if m.frameCount%2 == 0 {
			m.trail = append(m.trail, [2]int{m.cursorX, m.cursorY})
			if len(m.trail) > m.maxTrail {
				m.trail = m.trail[1:]
			}
		}

		// Maneja los clics
		if msg.Type == tea.MouseLeft {
			m.clickActive = true
			m.clickX, m.clickY = msg.X, msg.Y
			m.clickRadius = 1
			m.clickStyle = neonColors[m.frameCount%len(neonColors)]
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tickMsg:
		m.frameCount++

		// Actualiza las partículas con harmonica
		for i := range m.particles {
			// Calcula la posición objetivo con offset
			angle := float64(i) * (2 * math.Pi / float64(len(m.particles)))
			offset := 5.0 + 1.5*math.Sin(float64(m.frameCount)/10.0)

			targetX := float64(m.cursorX) + math.Cos(angle)*offset
			targetY := float64(m.cursorY) + math.Sin(angle)*offset

			m.particles[i].targetX = targetX
			m.particles[i].targetY = targetY

			// Actualiza las posiciones con harmonica
			m.particles[i].x, m.particles[i].xVel = m.particles[i].springX.Update(
				m.particles[i].x, m.particles[i].xVel, m.particles[i].targetX)
			m.particles[i].y, m.particles[i].yVel = m.particles[i].springY.Update(
				m.particles[i].y, m.particles[i].yVel, m.particles[i].targetY)
		}

		// Actualiza el efecto de clic
		if m.clickActive {
			m.clickRadius++
			if m.clickRadius > m.clickMaxRadius {
				m.clickActive = false
			}
		}

		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	// Creamos una matriz de caracteres para representar la pantalla
	screen := make([][]string, m.height)
	for y := range screen {
		screen[y] = make([]string, m.width)
		for x := range screen[y] {
			screen[y][x] = " "
		}
	}

	// Dibuja el rastro
	for i, pos := range m.trail {
		if pos[0] >= 0 && pos[0] < m.width && pos[1] >= 0 && pos[1] < m.height {
			opacity := float64(i) / float64(len(m.trail))
			idx := int(opacity * float64(len(neonColors)))
			if idx >= len(neonColors) {
				idx = len(neonColors) - 1
			}
			screen[pos[1]][pos[0]] = neonColors[idx].Render("·")
		}
	}

	// Dibuja el efecto de clic
	if m.clickActive {
		for y := m.clickY - m.clickRadius; y <= m.clickY + m.clickRadius; y++ {
			for x := m.clickX - m.clickRadius; x <= m.clickX + m.clickRadius; x++ {
				// Comprueba si el punto está dentro del círculo y de la pantalla
				distance := math.Sqrt(math.Pow(float64(x-m.clickX), 2) + math.Pow(float64(y-m.clickY), 2))
				if distance <= float64(m.clickRadius) && distance > float64(m.clickRadius-1) &&
					 x >= 0 && x < m.width && y >= 0 && y < m.height {
					screen[y][x] = m.clickStyle.Render("○")
				}
			}
		}
	}

	// Dibuja las partículas
	for i, p := range m.particles {
		x, y := int(p.x+0.5), int(p.y+0.5)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			// Alterna entre caracteres para crear variación
			char := "●"
			if i%2 == 0 {
				char = "◆"
			} else if i%3 == 0 {
				char = "■"
			}
			screen[y][x] = p.style.Render(char)
		}
	}

	// Dibuja el cursor
	if m.cursorX >= 0 && m.cursorX < m.width && m.cursorY >= 0 && m.cursorY < m.height {
		screen[m.cursorY][m.cursorX] = neonPink.Render("█")
	}

	// Convierte la matriz en un solo string
	var result strings.Builder
	for _, row := range screen {
		result.WriteString(strings.Join(row, ""))
		result.WriteString("\n")
	}

	// Agrega instrucciones
	result.WriteString("\n" + neonBlue.Render("Mueve el mouse - Haz clic para efectos - q para salir"))

	return result.String()
}

func main() {
	p := tea.NewProgram(initialModel(), 
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}