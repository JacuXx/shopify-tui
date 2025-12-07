package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Configuracion struct {
	Tiendas []Tienda `json:"tiendas"`
}

func obtenerDirectorioBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "shopify-tui"), nil
}

func obtenerRutaConfig() (string, error) {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores.json"), nil
}

func obtenerDirectorioStores() (string, error) {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirBase, "stores"), nil
}

func crearDirectorioBase() error {
	dirBase, err := obtenerDirectorioBase()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirBase, 0755)
}

func crearDirectorioStores() error {
	dirStores, err := obtenerDirectorioStores()
	if err != nil {
		return err
	}
	return os.MkdirAll(dirStores, 0755)
}

func crearDirectorioTienda(nombreTienda string) (string, error) {

	if err := crearDirectorioStores(); err != nil {
		return "", err
	}

	nombreCarpeta := sanitizarNombre(nombreTienda)

	dirStores, err := obtenerDirectorioStores()
	if err != nil {
		return "", err
	}

	rutaTienda := filepath.Join(dirStores, nombreCarpeta)

	if err := os.MkdirAll(rutaTienda, 0755); err != nil {
		return "", err
	}

	return rutaTienda, nil
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

	nombre = strings.Trim(nombre, "-")

	return nombre
}

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

func eliminarTienda(tiendas []Tienda, indice int) []Tienda {
	if indice < 0 || indice >= len(tiendas) {
		return tiendas
	}
	return append(tiendas[:indice], tiendas[indice+1:]...)
}

func existeDirectorio(ruta string) bool {
	info, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
