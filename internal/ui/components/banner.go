package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func lerpColor(c1, c2 colorful.Color, t float64) colorful.Color {
	return c1.BlendHcl(c2, t).Clamped()
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%02x%02x%02x",
		int(c.R*255), int(c.G*255), int(c.B*255))
}

var pixelFont = map[rune][]string{
	'S': {"01110", "10001", "10000", "01110", "00001", "10001", "01110"},
	'H': {"10001", "10001", "10001", "11111", "10001", "10001", "10001"},
	'O': {"01110", "10001", "10001", "10001", "10001", "10001", "01110"},
	'P': {"11110", "10001", "10001", "11110", "10000", "10000", "10000"},
	'I': {"11111", "00100", "00100", "00100", "00100", "00100", "11111"},
	'F': {"11111", "10000", "10000", "11110", "10000", "10000", "10000"},
	'Y': {"10001", "10001", "01010", "00100", "00100", "00100", "00100"},
	'T': {"11111", "00100", "00100", "00100", "00100", "00100", "00100"},
	'U': {"10001", "10001", "10001", "10001", "10001", "10001", "01110"},
	' ': {"000", "000", "000", "000", "000", "000", "000"},
}

const charHeight = 7

func buildGrid(text []rune) (grid [][]bool, totalCols int) {
	for i, ch := range text {
		if g, ok := pixelFont[ch]; ok {
			totalCols += len(g[0])
		}
		if i < len(text)-1 {
			totalCols++
		}
	}

	grid = make([][]bool, charHeight)
	for r := range grid {
		grid[r] = make([]bool, totalCols)
	}

	col := 0
	for ci, ch := range text {
		g, ok := pixelFont[ch]
		if !ok {
			continue
		}
		for row := range charHeight {
			for pc, px := range g[row] {
				grid[row][col+pc] = px == '1'
			}
		}
		col += len(g[0])
		if ci < len(text)-1 {
			col++
		}
	}
	return
}

func RenderBanner() string {
	text := []rune("SHOPIFY TUI")

	startColor, _ := colorful.Hex("#5E8E3E")
	endColor, _ := colorful.Hex("#96BF48")
	shadowColor := lipgloss.Color("#1a3a10")

	grid, totalCols := buildGrid(text)

	renderRows := charHeight + 1
	renderCols := totalCols + 1

	var banner strings.Builder

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#96BF48")).
		Bold(true)
	banner.WriteString(promptStyle.Render(">") + "\n\n")

	for row := range renderRows {
		for col := range renderCols {
			mainOn := row < charHeight && col < totalCols && grid[row][col]

			shadowOn := row > 0 && col > 0 &&
				(row-1) < charHeight && (col-1) < totalCols &&
				grid[row-1][col-1]

			switch {
			case mainOn:
				t := float64(col) / float64(totalCols-1)
				c := lerpColor(startColor, endColor, t)
				style := lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorToHex(c))).
					Bold(true)
				banner.WriteString(style.Render("██"))
			case shadowOn:
				style := lipgloss.NewStyle().Foreground(shadowColor)
				banner.WriteString(style.Render("░░"))
			default:
				banner.WriteString("  ")
			}
		}
		banner.WriteString("\n")
	}

	banner.WriteString("\n")

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	banner.WriteString(subtitleStyle.Render("  Shopify store manager — terminal edition") + "\n")

	return banner.String()
}

var _ = strings.Builder{}
