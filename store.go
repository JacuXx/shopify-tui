// store.go - Manejo de configuración y persistencia
// Este archivo se encarga de guardar y cargar las tiendas desde un archivo JSON
// También maneja la creación de directorios para cada tienda

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Configuracion representa la estructura del archivo JSON
type Configuracion struct {
	Tiendas []Tienda `json:"tiendas"`
}

// obtenerDirectorioBase retorna el directorio base donde se guardan las tiendas
// Por defecto: ~/.config/shopify-tui/
func obtenerDirectorioBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "shopify-tui"), nil
}

// obtenerRutaConfig retorna la ruta completa del archivo de configuración
// ~/.config/shopify-tui/stores.json
func obtenerRutaConfig() (string, error) {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores.json"), nil
}

// obtenerDirectorioStores retorna el directorio donde se guardan los temas
// ~/.config/shopify-tui/stores/
func obtenerDirectorioStores() (string, error) {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores"), nil
}

// crearDirectorioBase crea el directorio de configuración si no existe
func crearDirectorioBase() error {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirBase, 0755)
}

// crearDirectorioStores crea el directorio stores/ si no existe
func crearDirectorioStores() error {
	dirStores, err := obtenerDirectorioStores()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirStores, 0755)
}

// crearDirectorioTienda crea un directorio para una tienda específica
// Retorna la ruta completa del directorio creado
// Ejemplo: ~/.config/shopify-tui/stores/mi-tienda/
func crearDirectorioTienda(nombreTienda string) (string, error) {
	// Crear directorio stores/ si no existe
	if err := crearDirectorioStores(); err != nil {
		return "", err
	}

	// Sanitizar el nombre para usarlo como nombre de carpeta
	nombreCarpeta := sanitizarNombre(nombreTienda)

	dirStores, err := obtenerDirectorioStores()
	if err != nil {
		return "", err
	}

	rutaTienda := filepath.Join(dirStores, nombreCarpeta)

	// Crear el directorio de la tienda
	if err := os.MkdirAll(rutaTienda, 0755); err != nil {
		return "", err
	}

	return rutaTienda, nil
}

// sanitizarNombre convierte un nombre en algo seguro para usar como carpeta
// Ejemplo: "Mi Tienda Principal" -> "mi-tienda-principal"
func sanitizarNombre(nombre string) string {
	// Convertir a minúsculas
	nombre = strings.ToLower(nombre)

	// Reemplazar espacios y caracteres especiales con guiones
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

	// Eliminar guiones duplicados
	for strings.Contains(nombre, "--") {
		nombre = strings.ReplaceAll(nombre, "--", "-")
	}

	// Eliminar guiones al inicio y final
	nombre = strings.Trim(nombre, "-")

	return nombre
}

// cargarTiendas lee las tiendas desde el archivo JSON
func cargarTiendas() ([]Tienda, error) {
	rutaArchivo, err := obtenerRutaConfig()
	if err != nil {
		return nil, err
	}

	datos, err := os.ReadFile(rutaArchivo)
	if err != nil {
		if os.IsNotExist(err) {
			return []Tienda{}, nil
		}
		return nil, err
	}

	var config Configuracion
	if err := json.Unmarshal(datos, &config); err != nil {
		return nil, err
	}

	return config.Tiendas, nil
}

// guardarTiendas guarda las tiendas en el archivo JSON
func guardarTiendas(tiendas []Tienda) error {
	if err := crearDirectorioBase(); err != nil {
		return err
	}

	rutaArchivo, err := obtenerRutaConfig()
	if err != nil {
		return err
	}

	config := Configuracion{
		Tiendas: tiendas,
	}

	datos, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rutaArchivo, datos, 0644)
}

// eliminarTienda elimina una tienda por su índice
func eliminarTienda(tiendas []Tienda, indice int) []Tienda {
	if indice < 0 || indice >= len(tiendas) {
		return tiendas
	}
	return append(tiendas[:indice], tiendas[indice+1:]...)
}

// existeDirectorio verifica si un directorio existe
func existeDirectorio(ruta string) bool {
	info, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
