// store.go - Manejo de configuración y persistencia
// Este archivo se encarga de guardar y cargar las tiendas desde un archivo JSON
// La configuración se guarda en ~/.config/shopify-tui/stores.json

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Configuracion representa la estructura del archivo JSON
type Configuracion struct {
	Tiendas []Tienda `json:"tiendas"` // Lista de tiendas guardadas
}

// obtenerRutaConfig retorna la ruta completa del archivo de configuración
// Ejemplo: /home/usuario/.config/shopify-tui/stores.json
func obtenerRutaConfig() (string, error) {
	// os.UserHomeDir() obtiene el directorio home del usuario actual
	// En Linux: /home/usuario
	// En macOS: /Users/usuario
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// filepath.Join une las partes de la ruta de forma segura
	// (usa / en Linux/Mac y \ en Windows)
	rutaArchivo := filepath.Join(home, ".config", "shopify-tui", "stores.json")

	return rutaArchivo, nil
}

// crearDirectorioConfig crea el directorio de configuración si no existe
// Equivalente a: mkdir -p ~/.config/shopify-tui
func crearDirectorioConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rutaDir := filepath.Join(home, ".config", "shopify-tui")

	// os.MkdirAll crea todos los directorios necesarios
	// 0755 son los permisos: rwxr-xr-x (lectura/escritura para el dueño)
	return os.MkdirAll(rutaDir, 0755)
}

// cargarTiendas lee las tiendas desde el archivo JSON
// Si el archivo no existe, retorna una lista vacía (no es un error)
func cargarTiendas() ([]Tienda, error) {
	rutaArchivo, err := obtenerRutaConfig()
	if err != nil {
		return nil, err
	}

	// Leer todo el contenido del archivo
	datos, err := os.ReadFile(rutaArchivo)
	if err != nil {
		// Si el archivo no existe, no es un error - solo retornamos lista vacía
		if os.IsNotExist(err) {
			return []Tienda{}, nil
		}
		// Si es otro tipo de error, lo retornamos
		return nil, err
	}

	// Parsear el JSON a nuestra estructura
	var config Configuracion
	if err := json.Unmarshal(datos, &config); err != nil {
		return nil, err
	}

	return config.Tiendas, nil
}

// guardarTiendas guarda las tiendas en el archivo JSON
// Crea el directorio si no existe
func guardarTiendas(tiendas []Tienda) error {
	// Primero, asegurarnos de que el directorio existe
	if err := crearDirectorioConfig(); err != nil {
		return err
	}

	rutaArchivo, err := obtenerRutaConfig()
	if err != nil {
		return err
	}

	// Crear la estructura de configuración
	config := Configuracion{
		Tiendas: tiendas,
	}

	// json.MarshalIndent convierte la estructura a JSON con formato legible
	// "" es el prefijo (vacío)
	// "  " es la indentación (2 espacios)
	datos, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Escribir el archivo
	// 0644 son los permisos: rw-r--r-- (lectura/escritura para dueño, solo lectura para otros)
	return os.WriteFile(rutaArchivo, datos, 0644)
}

// eliminarTienda elimina una tienda por su índice
// Retorna la lista actualizada
func eliminarTienda(tiendas []Tienda, indice int) []Tienda {
	// Verificar que el índice sea válido
	if indice < 0 || indice >= len(tiendas) {
		return tiendas
	}

	// En Go, para eliminar un elemento de un slice:
	// append(slice[:indice], slice[indice+1:]...)
	// Esto une todo lo que está antes del índice con todo lo que está después
	return append(tiendas[:indice], tiendas[indice+1:]...)
}
