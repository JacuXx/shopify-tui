package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
)

func RenderBanner() string {
	myFigure := figure.NewFigure("SHOPIFY TUI", "slant", true)
	asciiArt := myFigure.String()
	verde := lipgloss.NewStyle().Foreground(lipgloss.Color("#95BF47")).Bold(true)
	lineas := strings.Split(asciiArt, "\n")
	var banner strings.Builder
	for _, linea := range lineas {
		if strings.TrimSpace(linea) == "" {
			continue
		}
		banner.WriteString(verde.Render(linea))
		banner.WriteString("\n")
	}
	return banner.String()
}
