package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	slideWidth  = 40
	slideHeight = 15
)

type model struct {
	slides     []slide
	currentIdx int
}

type slide interface {
	View() string
	Update(msg tea.Msg) (slide, tea.Cmd)
	Init() tea.Cmd
}

type creditsSlide struct {
	credits     []string
	currentLine int
	showAll     bool
}

type contentSlide struct {
	title string
	body  string
}

type barChartSlide struct {
	values      []int
	targets     []int
	labels      []string
	maxValue    int
	animating   bool
	initialized bool
}

type particleSlide struct {
	particles   []particle
	initialized bool
}

type particle struct {
	x, y        float64
	vx, vy      float64
	char        rune
	lifespan    int
	currentLife int
}

type gradientSlide struct {
	text         string
	startColor   lipgloss.Color
	endColor     lipgloss.Color
	progress     float64
	direction    int
	initialized  bool
}

type tickMsg struct{}

func initialModel() model {
	creditLines := []string{
		"Starring",
		"dev1 Agustín",
		"",
		"Lead Architect",
		"dev1 Agustín",
		"",
		"Backend Developer",
		"Jane Doe",
		"",
		"Frontend Developer",
		"John Smith",
		"",
		"Database Administrator",
		"Alex Johnson",
		"",
		"DevOps Engineer",
		"Maria García",
		"",
		"Project Manager",
		"Chris Williams",
		"",
		"A Charm Production",
		"2025",
	}

	barChartData := []int{0, 0, 0, 0, 0}
	barChartTargets := []int{22, 16, 31, 18, 27}
	barChartLabels := []string{"Proyecto A", "Proyecto B", "Proyecto C", "Proyecto D", "Proyecto E"}

	slides := []slide{
		&creditsSlide{
			credits:     creditLines,
			currentLine: -5,
		},
		&contentSlide{
			title: "Navegación de Slides",
			body:  "Este es un proyecto demostrativo de una CLI con slides.\n\nUsa las flechas ← → para navegar entre slides.\n\nPresiona 'q' para salir.",
		},
		&barChartSlide{
			values:      barChartData,
			targets:     barChartTargets,
			labels:      barChartLabels,
			maxValue:    35,
			animating:   true,
			initialized: false,
		},
		&particleSlide{
			particles:   make([]particle, 0),
			initialized: false,
		},
		&gradientSlide{
			text:        "Este texto cambiará de color gradualmente",
			startColor:  lipgloss.Color("#FF0000"),
			endColor:    lipgloss.Color("#0000FF"),
			progress:    0.0,
			direction:   1,
			initialized: false,
		},
	}

	return model{
		slides:     slides,
		currentIdx: 0,
	}
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			Width(slideWidth - 4).
			Align(lipgloss.Center)

	slideStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			Width(slideWidth).
			Height(slideHeight)

	bodyStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			Width(slideWidth - 4)

	creditStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Align(lipgloss.Center).
			Width(slideWidth)

	creditTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Align(lipgloss.Center).
			Width(slideWidth)

	barStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#7D56F4"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Width(12).
			Align(lipgloss.Left)

	particleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00"))
)

func (c *creditsSlide) Init() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (c *creditsSlide) Update(msg tea.Msg) (slide, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		if c.showAll {
			return c, nil
		}

		c.currentLine++

		if c.currentLine > len(c.credits)+slideHeight {
			c.showAll = true
			return c, nil
		}

		return c, tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return c, nil
}

func (c *creditsSlide) View() string {
	if c.showAll {
		var sb strings.Builder
		for i, line := range c.credits {
			if i%2 == 0 {
				sb.WriteString(creditTitleStyle.Render(line))
			} else {
				sb.WriteString(creditStyle.Render(line))
			}
			sb.WriteString("\n")
		}
		return slideStyle.Render(sb.String())
	}

	var sb strings.Builder
	for i := 0; i < slideHeight; i++ {
		lineIdx := c.currentLine - slideHeight + i
		if lineIdx >= 0 && lineIdx < len(c.credits) {
			if lineIdx%2 == 0 {
				sb.WriteString(creditTitleStyle.Render(c.credits[lineIdx]))
			} else {
				sb.WriteString(creditStyle.Render(c.credits[lineIdx]))
			}
		}
		sb.WriteString("\n")
	}

	return slideStyle.Render(sb.String())
}

func (s *contentSlide) Init() tea.Cmd {
	return nil
}

func (s *contentSlide) Update(msg tea.Msg) (slide, tea.Cmd) {
	return s, nil
}

func (s *contentSlide) View() string {
	return slideStyle.Render(
		fmt.Sprintf("%s\n%s",
			titleStyle.Render(s.title),
			bodyStyle.Render(s.body),
		),
	)
}

func (b *barChartSlide) Init() tea.Cmd {
	if !b.initialized {
		b.initialized = true
		return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return nil
}

func (b *barChartSlide) Update(msg tea.Msg) (slide, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		if !b.animating {
			return b, nil
		}

		allComplete := true
		for i := range b.values {
			if b.values[i] < b.targets[i] {
				b.values[i] += 1 + rand.Intn(2)
				if b.values[i] > b.targets[i] {
					b.values[i] = b.targets[i]
				} else {
					allComplete = false
				}
			}
		}

		if allComplete {
			b.animating = false
			return b, nil
		}

		return b, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return b, nil
}

func (b *barChartSlide) View() string {
	chartHeight := 10
	barWidth := 8
	maxChartValue := b.maxValue

	title := titleStyle.Render("Rendimiento por Proyecto")
	var sb strings.Builder
	sb.WriteString(title + "\n\n")

	for i, value := range b.values {
		barLength := int(float64(value) / float64(maxChartValue) * float64(chartHeight))
		label := b.labels[i]

		sb.WriteString(labelStyle.Render(label) + " ")

		barText := fmt.Sprintf("%d", value)
		bar := strings.Repeat(" ", barWidth)
		barRendered := barStyle.Copy().Width(barWidth).Render(bar)

		for row := 0; row < chartHeight; row++ {
			y := chartHeight - row - 1
			if y < barLength {
				if row == chartHeight - barLength {
					valueStyle := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#000000")).
						Background(lipgloss.Color("#7D56F4")).
						Width(barWidth).
						Align(lipgloss.Center)
					sb.WriteString(valueStyle.Render(barText))
				} else {
					sb.WriteString(barRendered)
				}
			} else {
				sb.WriteString(strings.Repeat(" ", barWidth))
			}
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}

	return slideStyle.Render(sb.String())
}

func (p *particleSlide) Init() tea.Cmd {
	if !p.initialized {
		p.initialized = true
		p.particles = make([]particle, 0)
		return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return nil
}

func (p *particleSlide) Update(msg tea.Msg) (slide, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		if rand.Intn(3) == 0 {
			if len(p.particles) < 100 {
				chars := []rune{'*', '+', '.', '·', '•', '°', '✧', '✦', '✴', '✹'}
				x := float64(slideWidth/2) + rand.Float64()*4 - 2
				y := float64(slideHeight/2)
				angle := rand.Float64() * 2 * math.Pi
				speed := 0.2 + rand.Float64()*0.4

				p.particles = append(p.particles, particle{
					x:           x,
					y:           y,
					vx:          math.Cos(angle) * speed,
					vy:          math.Sin(angle) * speed,
					char:        chars[rand.Intn(len(chars))],
					lifespan:    50 + rand.Intn(50),
					currentLife: 0,
				})
			}
		}

		for i := 0; i < len(p.particles); i++ {
			p.particles[i].x += p.particles[i].vx
			p.particles[i].y += p.particles[i].vy
			p.particles[i].currentLife++

			if p.particles[i].currentLife >= p.particles[i].lifespan ||
				 p.particles[i].x < 0 || p.particles[i].x >= float64(slideWidth) ||
				 p.particles[i].y < 0 || p.particles[i].y >= float64(slideHeight) {
				p.particles = append(p.particles[:i], p.particles[i+1:]...)
				i--
			}
		}

		return p, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return p, nil
}

func (p *particleSlide) View() string {
	grid := make([][]rune, slideHeight)
	for i := range grid {
		grid[i] = make([]rune, slideWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for _, particle := range p.particles {
		x, y := int(particle.x), int(particle.y)
		if x >= 0 && x < slideWidth && y >= 0 && y < slideHeight {
			grid[y][x] = particle.char
		}
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Simulación de Partículas") + "\n\n")

	for _, row := range grid {
		sb.WriteString(particleStyle.Render(string(row)) + "\n")
	}

	return slideStyle.Render(sb.String())
}

func (g *gradientSlide) Init() tea.Cmd {
	if !g.initialized {
		g.initialized = true
		return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return nil
}

func (g *gradientSlide) Update(msg tea.Msg) (slide, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		g.progress += 0.02 * float64(g.direction)

		if g.progress >= 1.0 {
			g.progress = 1.0
			g.direction = -1
		} else if g.progress <= 0.0 {
			g.progress = 0.0
			g.direction = 1
		}

		return g, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return g, nil
}

func (g *gradientSlide) View() string {
	startRed, startGreen, startBlue := parseHexColor(string(g.startColor))
	endRed, endGreen, endBlue := parseHexColor(string(g.endColor))

	currentRed := int(float64(startRed) + g.progress*float64(endRed-startRed))
	currentGreen := int(float64(startGreen) + g.progress*float64(endGreen-startGreen))
	currentBlue := int(float64(startBlue) + g.progress*float64(endBlue-startBlue))

	currentColor := fmt.Sprintf("#%02X%02X%02X", currentRed, currentGreen, currentBlue)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentColor)).
		Bold(true).
		Width(slideWidth - 4).
		Align(lipgloss.Center)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Cambios de Estilo Progresivos") + "\n\n")
	sb.WriteString("\n\n\n")
	sb.WriteString(textStyle.Render(g.text) + "\n\n")
	sb.WriteString(textStyle.Render(fmt.Sprintf("Color actual: %s (Progreso: %.0f%%)", currentColor, g.progress*100)))

	return slideStyle.Render(sb.String())
}

func parseHexColor(hexColor string) (int, int, int) {
	if hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	var r, g, b int
	fmt.Sscanf(hexColor, "%02X%02X%02X", &r, &g, &b)
	return r, g, b
}

func (m model) Init() tea.Cmd {
	return m.slides[m.currentIdx].Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "right", "l", "n", " ":
			if m.currentIdx < len(m.slides)-1 {
				m.currentIdx++
				return m, m.slides[m.currentIdx].Init()
			}

		case "left", "h", "p":
			if m.currentIdx > 0 {
				m.currentIdx--
				return m, m.slides[m.currentIdx].Init()
			}
		}
	}

	var cmd tea.Cmd
	m.slides[m.currentIdx], cmd = m.slides[m.currentIdx].Update(msg)
	return m, cmd
}

func (m model) View() string {
	slideView := m.slides[m.currentIdx].View()

	nav := fmt.Sprintf("\n[%d/%d] Use ← → para navegar, 'q' para salir", 
		m.currentIdx+1, len(m.slides))

	return slideView + "\n" + nav
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}