package components

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//go:embed shopify-logo.png
var shopifyLogoData []byte

func lerpColor(c1, c2 colorful.Color, t float64) colorful.Color {
	return c1.BlendHcl(c2, t).Clamped()
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%02x%02x%02x",
		int(c.R*255), int(c.G*255), int(c.B*255))
}

func resizeNearest(img image.Image, w, h int) image.Image {
	src := img.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			sx := x*src.Dx()/w + src.Min.X
			sy := y*src.Dy()/h + src.Min.Y
			dst.Set(x, y, img.At(sx, sy))
		}
	}
	return dst
}

func renderLogo(data []byte, termW, termH int) []string {
	rows := make([]string, termH)
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return rows
	}
	resized := resizeNearest(img, termW, termH*2)
	for row := range termH {
		var sb strings.Builder
		for col := range termW {
			_, _, _, topA := resized.At(col, row*2).RGBA()
			_, _, _, botA := resized.At(col, row*2+1).RGBA()
			tR, tG, tB, _ := resized.At(col, row*2).RGBA()
			bR, bG, bB, _ := resized.At(col, row*2+1).RGBA()
			topHex := fmt.Sprintf("#%02x%02x%02x", tR>>8, tG>>8, tB>>8)
			botHex := fmt.Sprintf("#%02x%02x%02x", bR>>8, bG>>8, bB>>8)
			topSolid := topA > 0x7fff
			botSolid := botA > 0x7fff
			switch {
			case topSolid && botSolid:
				sb.WriteString(lipgloss.NewStyle().
					Background(lipgloss.Color(topHex)).
					Foreground(lipgloss.Color(botHex)).
					Render("▄"))
			case topSolid:
				sb.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(topHex)).
					Render("▀"))
			case botSolid:
				sb.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(botHex)).
					Render("▄"))
			default:
				sb.WriteString(" ")
			}
		}
		rows[row] = sb.String()
	}
	return rows
}

var pixelFont = map[rune][]string{
	'S': {"01110", "10000", "01110", "00001", "01110"},
	'H': {"10001", "10001", "11111", "10001", "10001"},
	'O': {"01110", "10001", "10001", "10001", "01110"},
	'P': {"11110", "10001", "11110", "10000", "10000"},
	'I': {"11111", "00100", "00100", "00100", "11111"},
	'F': {"11111", "10000", "11110", "10000", "10000"},
	'Y': {"10001", "01010", "00100", "00100", "00100"},
	'T': {"11111", "00100", "00100", "00100", "00100"},
	'U': {"10001", "10001", "10001", "10001", "01110"},
	' ': {"000", "000", "000", "000", "000"},
}

const charHeight = 5

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

func BannerHeight(ancho int) int {
	switch {
	case ancho > 0 && ancho < 64:
		return 3
	case ancho > 0 && ancho < 84:
		return 9
	default:
		return 12
	}
}

func RenderBanner(ancho int) string {
	if ancho > 0 && ancho < 64 {
		return renderBannerMinimal()
	}
	if ancho > 0 && ancho < 84 {
		return renderBannerCompacto()
	}
	return renderBannerCompleto()
}

func renderBannerMinimal() string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#96BF48")).
		Bold(true)
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	b.WriteString(titleStyle.Render("SHOPIFY TUI") + "\n")
	b.WriteString(subtitleStyle.Render("store manager — terminal edition") + "\n")
	return b.String()
}

func renderBannerCompacto() string {
	text := []rune("SHOPIFY TUI")

	startColor, _ := colorful.Hex("#5E8E3E")
	endColor, _ := colorful.Hex("#96BF48")
	shadowColor := lipgloss.Color("#1a3a10")

	grid, totalCols := buildGrid(text)
	renderRows := charHeight + 1
	renderCols := totalCols + 1

	var banner strings.Builder
	banner.WriteString("\n")

	for row := range renderRows {
		for col := range renderCols {
			mainOn := row < charHeight && col < totalCols && grid[row][col]
			shadowOn := row > 0 && col > 0 &&
				(row-1) < charHeight && (col-1) < totalCols &&
				grid[row-1][col-1]

			t := float64(col) / float64(totalCols-1)
			c := lerpColor(startColor, endColor, t)
			bold := lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))).Bold(true)

			switch {
			case mainOn:
				banner.WriteString(bold.Render("█"))
			case shadowOn:
				banner.WriteString(lipgloss.NewStyle().Foreground(shadowColor).Render("░"))
			default:
				banner.WriteString(" ")
			}
		}
		banner.WriteString("\n")
	}

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	banner.WriteString(subtitleStyle.Render("  Shopify store manager — terminal edition") + "\n")

	return banner.String()
}

func renderBannerCompleto() string {
	text := []rune("SHOPIFY TUI")

	startColor, _ := colorful.Hex("#5E8E3E")
	endColor, _ := colorful.Hex("#96BF48")
	shadowColor := lipgloss.Color("#1a3a10")

	grid, totalCols := buildGrid(text)

	const logoTermRows = 8
	const logoTermW = logoTermRows * 2
	logoRows := renderLogo(shopifyLogoData, logoTermW, logoTermRows)

	renderRows := charHeight + 1
	renderCols := totalCols + 1
	textStartRow := logoTermRows - renderRows

	var banner strings.Builder
	banner.WriteString("\n")

	for row := range logoTermRows {
		banner.WriteString(logoRows[row])
		banner.WriteString("  ")

		textRow := row - textStartRow
		if textRow >= 0 && textRow < renderRows {
			for col := range renderCols {
				mainOn := textRow < charHeight && col < totalCols && grid[textRow][col]
				shadowOn := textRow > 0 && col > 0 &&
					(textRow-1) < charHeight && (col-1) < totalCols &&
					grid[textRow-1][col-1]

				t := float64(col) / float64(totalCols-1)
				c := lerpColor(startColor, endColor, t)
				bold := lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))).Bold(true)

				switch {
				case mainOn:
					banner.WriteString(bold.Render("█"))
				case shadowOn:
					banner.WriteString(lipgloss.NewStyle().Foreground(shadowColor).Render("░"))
				default:
					banner.WriteString(" ")
				}
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
