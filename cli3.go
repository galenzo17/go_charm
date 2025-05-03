				package main

				import (
					"fmt"
					"math/rand"
					"os"
					"time"

					"github.com/charmbracelet/lipgloss"
				)

				const (
					width  = 60
					height = 10
				)

				// Estilos neón
				var (
					carStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
					obstacleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
					groundStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
					scoreStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true)
					gameOverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8800")).Bold(true)
				)

				type Game struct {
					screen     [][]rune
					carPos     int
					carHeight  int
					obstacles  []int
					obstacleTypes []string
					score      int
					gameOver   bool
				}

				func NewGame() *Game {
					screen := make([][]rune, height)
					for i := range screen {
						screen[i] = make([]rune, width)
						for j := range screen[i] {
							screen[i][j] = ' '
						}
					}

					return &Game{
						screen:     screen,
						carPos:     0,
						carHeight:  height - 3,
						obstacles:  []int{width},
						obstacleTypes: []string{"^", "#", "@"},
						score:      0,
						gameOver:   false,
					}
				}

				func (g *Game) Update() {
					if g.gameOver {
						return
					}

					// Aumentar puntuación
					g.score++

					// Actualizar posición del coche (salto)
					if g.carPos > 0 {
						g.carPos--
					}

					// Actualizar obstáculos
					for i := range g.obstacles {
						g.obstacles[i]--

						// Comprobar colisiones
						if g.obstacles[i] == 5 && g.carHeight + g.carPos >= height-3 {
							g.gameOver = true
						}

						// Eliminar obstáculos fuera de pantalla
						if g.obstacles[i] < 0 {
							g.obstacles = append(g.obstacles[:i], g.obstacles[i+1:]...)
							break
						}
					}

					// Agregar nuevos obstáculos
					if rand.Intn(15) == 0 {
						g.obstacles = append(g.obstacles, width)
					}
				}

				func (g *Game) Draw() {
					// Limpiar pantalla
					fmt.Print("\033[H\033[2J")

					// Limpiar buffer
					for i := range g.screen {
						for j := range g.screen[i] {
							g.screen[i][j] = ' '
						}
					}

					// Dibujar coche
					carY := height - 3 - g.carPos
					g.screen[carY][5] = '>'
					g.screen[carY][4] = '-'
					g.screen[carY+1][5] = 'O'
					g.screen[carY-1][5] = '^'

					// Dibujar obstáculos
					for _, obsX := range g.obstacles {
						if obsX >= 0 && obsX < width {
							obsType := g.obstacleTypes[rand.Intn(len(g.obstacleTypes))]
							g.screen[height-3][obsX] = []rune(obsType)[0]
						}
					}

					// Dibujar suelo
					for i := 0; i < width; i++ {
						g.screen[height-2][i] = '='
					}

					// Imprimir puntuación
					fmt.Println(scoreStyle.Render(fmt.Sprintf("SCORE: %d", g.score)))

					// Imprimir marco superior
					fmt.Print(groundStyle.Render("+"))
					for i := 0; i < width; i++ {
						fmt.Print(groundStyle.Render("-"))
					}
					fmt.Println(groundStyle.Render("+"))

					// Imprimir pantalla
					for i := 0; i < height; i++ {
						fmt.Print(groundStyle.Render("|"))
						for j := 0; j < width; j++ {
							char := string(g.screen[i][j])
							if char == ">" || char == "-" || char == "O" || char == "^" {
								fmt.Print(carStyle.Render(char))
							} else if char == "^" || char == "#" || char == "@" {
								fmt.Print(obstacleStyle.Render(char))
							} else if char == "=" {
								fmt.Print(groundStyle.Render(char))
							} else {
								fmt.Print(char)
							}
						}
						fmt.Println(groundStyle.Render("|"))
					}

					// Imprimir marco inferior
					fmt.Print(groundStyle.Render("+"))
					for i := 0; i < width; i++ {
						fmt.Print(groundStyle.Render("-"))
					}
					fmt.Println(groundStyle.Render("+"))

					// Instrucciones o mensaje de game over
					if g.gameOver {
						fmt.Println(gameOverStyle.Render("GAME OVER! Presiona ENTER para reiniciar"))
					} else {
						fmt.Println("Presiona ENTER para saltar, escribe 'q' para salir")
					}
				}

				func (g *Game) Jump() {
					if g.carPos == 0 {
						g.carPos = 4 // Altura del salto
					}
				}

				func main() {
					rand.Seed(time.Now().UnixNano())

					game := NewGame()

					fmt.Println(carStyle.Render("=== NEON CAR RUNNER ==="))
					fmt.Println(groundStyle.Render("Instrucciones:"))
					fmt.Println(obstacleStyle.Render("- Presiona ENTER para saltar"))
					fmt.Println(obstacleStyle.Render("- Escribe 'q' y presiona ENTER para salir"))
					fmt.Println(carStyle.Render("Presiona ENTER para comenzar..."))
					fmt.Scanln()

					ticker := time.NewTicker(100 * time.Millisecond)

					go func() {
						for range ticker.C {
							if !game.gameOver {
								game.Update()
								game.Draw()
							}
						}
					}()

					for {
						var input string
						fmt.Scanln(&input)

						if input == "q" {
							break
						}

						if game.gameOver {
							game = NewGame()
						} else {
							game.Jump()
						}
					}

					ticker.Stop()
					fmt.Println("¡Gracias por jugar!")
				}
				if screenY >= 0 && screenY < len(screen) && screenX >= 0 && screenX < len(screen[0]) {
					screen[screenY][screenX] = string(o.pattern[y][x])
				}
			}
		}
	}
}

// IsColliding checks if the obstacle is colliding with the car
func (o *Obstacle) IsColliding(car *Car) bool {
	// Simple rectangle collision
	return o.x < car.x+car.width && 
				 o.x+o.width > car.x &&
				 o.y < car.y+car.height && 
				 o.y+o.height > car.y
}

// IsOffScreen checks if the obstacle is off the screen
func (o *Obstacle) IsOffScreen() bool {
	return o.x+o.width < 0
}

// Game represents the game state
type Game struct {
	car        *Car
	obstacles  []*Obstacle
	screen     [][]string
	score      int
	speed      int
	nextObsPos int
	gameOver   bool
	keyPressed bool
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Initialize screen buffer
	screen := make([][]string, screenHeight)
	for i := range screen {
		screen[i] = make([]string, screenWidth)
		for j := range screen[i] {
			screen[i][j] = " "
		}
	}

	return &Game{
		car:        NewCar(),
		obstacles:  []*Obstacle{},
		screen:     screen,
		speed:      1,
		nextObsPos: screenWidth + rand.Intn(maxObsDist-minObsDist) + minObsDist,
	}
}

// HandleInput processes input
func (g *Game) HandleInput(ch byte) {
	if ch == ' ' || ch == 'w' || ch == 'W' {
		if g.gameOver {
			// Restart game
			*g = *NewGame()
		} else if g.car.y == carBaseY && !g.keyPressed {
			// Jump only if car is on the ground
			g.car.jump = true
			g.keyPressed = true
		}
	}

	if ch == 0 {
		g.keyPressed = false
	}
}

// Update updates the game state
func (g *Game) Update() {
	if g.gameOver {
		return
	}

	// Update car
	g.car.Update()

	// Update score and speed
	g.score++
	if g.score%500 == 0 {
		g.speed = min(g.score/500 + 1, 5) // Cap speed at 5
	}

	// Create new obstacle
	if g.nextObsPos <= screenWidth {
		g.obstacles = append(g.obstacles, NewObstacle(g.nextObsPos))
		g.nextObsPos = screenWidth + rand.Intn(maxObsDist-minObsDist) + minObsDist
	} else {
		g.nextObsPos -= g.speed
	}

	// Update obstacles
	for i := 0; i < len(g.obstacles); i++ {
		g.obstacles[i].Update(g.speed)

		// Check collision
		if g.obstacles[i].IsColliding(g.car) {
			g.gameOver = true
		}

		// Remove obstacles that are off-screen
		if g.obstacles[i].IsOffScreen() {
			g.obstacles = append(g.obstacles[:i], g.obstacles[i+1:]...)
			i--
		}
	}
}

// Draw renders the game
func (g *Game) Draw() {
	// Clear screen buffer
	for i := range g.screen {
		for j := range g.screen[i] {
			g.screen[i][j] = " "
		}
	}

	// Draw ground
	for i := 0; i < screenWidth; i++ {
		g.screen[groundY][i] = "-"
	}

	// Draw car
	g.car.Draw(g.screen)

	// Draw obstacles
	for _, obs := range g.obstacles {
		obs.Draw(g.screen)
	}

	// Clear terminal
	clearScreen()

	// Create the display with frames and styling
	output := ""

	// Draw score
	scoreText := fmt.Sprintf("SCORE: %d", g.score)
	output += scoreStyle.Render(scoreText) + "\n"

	// Draw top border
	border := "+"
	for i := 0; i < screenWidth; i++ {
		border += "-"
	}
	border += "+"
	output += groundStyle.Render(border) + "\n"

	// Draw game area with styled characters
	for _, row := range g.screen {
		line := "|"
		for _, char := range row {
			if char == " " {
				line += bgStyle.Render(" ")
			} else if char == "-" {
				line += groundStyle.Render(char)
			} else if char == "O" || char == "/" || char == "\\" || char == "_" || char == "-" {
				if g.gameOver {
					line += gameOverStyle.Render(char)
				} else {
					line += carStyle.Render(char)
				}
			} else if char == "^" || char == "|" || char == "/" || char == "\\" || char == "_" {
				line += obstStyle.Render(char)
			} else {
				line += groundStyle.Render(char)
			}
		}
		line += "|"
		output += groundStyle.Render(line) + "\n"
	}

	// Draw bottom border
	border = "+"
	for i := 0; i < screenWidth; i++ {
		border += "-"
	}
	border += "+"
	output += groundStyle.Render(border) + "\n"

	// Game over message
	if g.gameOver {
		gameOverMsg := "GAME OVER! PRESS SPACE TO RESTART"
		output += gameOverStyle.Render(gameOverMsg)
	} else {
		output += lipgloss.NewStyle().Foreground(neonOrange).Render("Press SPACE or W to jump")
	}

	fmt.Print(output)
}

// clearScreen clears the terminal
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Simple input loop (no special terminal handling required)
func simpleGameLoop(g *Game) {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			g.Update()
			g.Draw()
		}
	}()

	var input string
	for {
		fmt.Scanln(&input)
		if len(input) > 0 {
			if input == "q" || input == "Q" {
				clearScreen()
				return
			}
			g.HandleInput(' ') // Treat any input as space
		}
	}
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize the game
	game := NewGame()

	// Display welcome message
	clearScreen()
	fmt.Println(lipgloss.NewStyle().Foreground(neonPink).Bold(true).Render("=== NEON CAR RUNNER ==="))
	fmt.Println(lipgloss.NewStyle().Foreground(neonBlue).Render("Controls:"))
	fmt.Println(lipgloss.NewStyle().Foreground(neonGreen).Render("- Press ENTER to jump"))
	fmt.Println(lipgloss.NewStyle().Foreground(neonGreen).Render("- Type 'q' + ENTER to quit"))
	fmt.Println(lipgloss.NewStyle().Foreground(neonOrange).Render("Press ENTER to start..."))

	// Wait for key press
	fmt.Scanln()

	// Start the game loop
	simpleGameLoop(game)
}