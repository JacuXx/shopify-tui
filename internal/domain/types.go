package domain

// Vista representa el estado de navegación actual de la app.
type Vista int

const (
	VistaMenu Vista = iota
	VistaAgregarTienda
	VistaSeleccionarMetodo
	VistaInputGit
	VistaSeleccionarTienda
	VistaSeleccionarModo
	VistaLogs
	VistaServidores
	VistaPopup
)

// MetodoDescarga indica cómo se descarga el tema de una tienda.
type MetodoDescarga int

const (
	MetodoShopifyPull MetodoDescarga = iota
	MetodoGitClone
)

// Tienda representa una tienda Shopify registrada en la app.
type Tienda struct {
	Nombre string         `json:"nombre"`
	URL    string         `json:"url"`
	Ruta   string         `json:"ruta"`
	Metodo MetodoDescarga `json:"metodo"`
	GitURL string         `json:"git_url,omitempty"`
}
