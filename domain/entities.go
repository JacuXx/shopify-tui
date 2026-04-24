package domain

import (
	"time"
)

// Tienda representa una tienda Shopify registrada
type Tienda struct {
	ID     string
	Nombre string
	URL    string
	Ruta   string
	Metodo DownloadMethod
	GitURL string
}

// ServidorActivo representa un servidor de desarrollo en ejecución
type ServidorActivo struct {
	ID       string
	Tienda   Tienda
	Puerto   int
	Iniciado time.Time
	URL      string
	Activo   bool
	Logs     []string
}

// DownloadMethod representa el método para descargar temas
type DownloadMethod int

const (
	MethodShopifyPull DownloadMethod = iota
	MethodGitClone
)

func (m DownloadMethod) String() string {
	switch m {
	case MethodShopifyPull:
		return "shopify-pull"
	case MethodGitClone:
		return "git-clone"
	default:
		return "unknown"
	}
}
