package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/JacuXx/shopify-cli/internal/domain"
)

type configuracion struct {
	Tiendas []domain.Tienda `json:"tiendas"`
}

// JSONRepository implementa Repository usando un archivo JSON local.
// Path: ~/.config/shopify-tui/stores.json
type JSONRepository struct{}

// NewJSONRepository crea una instancia del repositorio JSON.
func NewJSONRepository() *JSONRepository {
	return &JSONRepository{}
}

func (r *JSONRepository) obtenerDirectorioBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "shopify-tui"), nil
}

func (r *JSONRepository) obtenerRutaConfig() (string, error) {
	dirBase, err := r.obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores.json"), nil
}

func (r *JSONRepository) obtenerDirectorioStores() (string, error) {
	dirBase, err := r.obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores"), nil
}

func (r *JSONRepository) crearDirectorioBase() error {
	dirBase, err := r.obtenerDirectorioBase()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirBase, 0755)
}

func (r *JSONRepository) crearDirectorioStores() error {
	dirStores, err := r.obtenerDirectorioStores()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirStores, 0755)
}

// CargarTiendas lee el archivo JSON y retorna la lista de tiendas.
func (r *JSONRepository) CargarTiendas() ([]domain.Tienda, error) {
	rutaArchivo, err := r.obtenerRutaConfig()
	if err != nil {
		return nil, err
	}

	datos, err := os.ReadFile(rutaArchivo)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Tienda{}, nil
		}
		return nil, err
	}

	var config configuracion
	if err := json.Unmarshal(datos, &config); err != nil {
		return nil, err
	}

	return config.Tiendas, nil
}

// GuardarTiendas escribe la lista de tiendas en el archivo JSON.
func (r *JSONRepository) GuardarTiendas(tiendas []domain.Tienda) error {
	if err := r.crearDirectorioBase(); err != nil {
		return err
	}

	rutaArchivo, err := r.obtenerRutaConfig()
	if err != nil {
		return err
	}

	config := configuracion{Tiendas: tiendas}

	datos, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rutaArchivo, datos, 0600)
}

// CrearDirectorioTienda crea el directorio local para una tienda y retorna su ruta.
func (r *JSONRepository) CrearDirectorioTienda(nombre string) (string, error) {
	if err := r.crearDirectorioStores(); err != nil {
		return "", err
	}

	nombreCarpeta := sanitizarNombre(nombre)

	dirStores, err := r.obtenerDirectorioStores()
	if err != nil {
		return "", err
	}

	rutaTienda := filepath.Join(dirStores, nombreCarpeta)

	if err := os.MkdirAll(rutaTienda, 0755); err != nil {
		return "", err
	}

	return rutaTienda, nil
}

// EliminarTienda elimina una tienda del slice por índice.
func (r *JSONRepository) EliminarTienda(tiendas []domain.Tienda, indice int) []domain.Tienda {
	if indice < 0 || indice >= len(tiendas) {
		return tiendas
	}
	return append(tiendas[:indice], tiendas[indice+1:]...)
}

// ExisteDirectorio verifica si una ruta existe y es un directorio.
func (r *JSONRepository) ExisteDirectorio(ruta string) bool {
	info, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && info.IsDir()
}

func sanitizarNombre(nombre string) string {
	nombre = strings.ToLower(nombre)

	replacer := strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	nombre = replacer.Replace(nombre)

	for strings.Contains(nombre, "--") {
		nombre = strings.ReplaceAll(nombre, "--", "-")
	}

	return strings.Trim(nombre, "-")
}
