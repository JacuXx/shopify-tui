// server.go - Manejo de servidores de desarrollo en background
// Este archivo contiene la lógica para ejecutar y gestionar múltiples
// servidores de desarrollo corriendo simultáneamente

package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// ServidorActivo representa un servidor theme dev corriendo en background
type ServidorActivo struct {
	Tienda    Tienda       // Tienda asociada al servidor
	Proceso   *exec.Cmd    // Proceso del comando
	Puerto    int          // Puerto donde corre (default 9292)
	Iniciado  time.Time    // Cuándo se inició
	URL       string       // URL local del servidor
	Activo    bool         // Si está corriendo
	Logs      []string     // Últimas líneas de logs
	LogsMutex sync.RWMutex // Mutex para los logs
	Stdin     io.WriteCloser // Para enviar input al proceso
}

// AgregarLog añade una línea de log (mantiene las últimas 100 líneas)
func (s *ServidorActivo) AgregarLog(linea string) {
	s.LogsMutex.Lock()
	defer s.LogsMutex.Unlock()
	
	s.Logs = append(s.Logs, linea)
	// Mantener solo las últimas 100 líneas
	if len(s.Logs) > 100 {
		s.Logs = s.Logs[len(s.Logs)-100:]
	}
}

// ObtenerLogs retorna una copia de los logs
func (s *ServidorActivo) ObtenerLogs() []string {
	s.LogsMutex.RLock()
	defer s.LogsMutex.RUnlock()
	
	copia := make([]string, len(s.Logs))
	copy(copia, s.Logs)
	return copia
}

// EnviarInput envía una tecla o texto al proceso
func (s *ServidorActivo) EnviarInput(input string) error {
	if s.Stdin == nil {
		return fmt.Errorf("stdin no disponible")
	}
	_, err := s.Stdin.Write([]byte(input))
	return err
}

// GestorServidores maneja todos los servidores activos
// Es thread-safe usando mutex
type GestorServidores struct {
	servidores map[string]*ServidorActivo // Mapa de servidores por nombre de tienda
	mutex      sync.RWMutex               // Para acceso concurrente seguro
	puertos    map[int]bool               // Puertos en uso
}

// gestorGlobal es la instancia global del gestor de servidores
var gestorGlobal = &GestorServidores{
	servidores: make(map[string]*ServidorActivo),
	puertos:    make(map[int]bool),
}

// ObtenerGestor retorna el gestor global de servidores
func ObtenerGestor() *GestorServidores {
	return gestorGlobal
}

// ObtenerPuertoDisponible encuentra el siguiente puerto disponible
// Empieza en 9292 (default de Shopify) y busca uno libre
func (g *GestorServidores) ObtenerPuertoDisponible() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	puerto := 9292
	for g.puertos[puerto] {
		puerto++
	}
	return puerto
}

// IniciarServidor inicia un nuevo servidor de desarrollo para una tienda
func (g *GestorServidores) IniciarServidor(tienda Tienda) (*ServidorActivo, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Verificar si ya hay un servidor para esta tienda
	if servidor, existe := g.servidores[tienda.Nombre]; existe && servidor.Activo {
		return nil, fmt.Errorf("ya hay un servidor activo para '%s'", tienda.Nombre)
	}

	// Obtener puerto disponible
	puerto := 9292
	for g.puertos[puerto] {
		puerto++
	}

	// Crear el comando
	cmd := exec.Command("shopify", "theme", "dev",
		"--store", tienda.URL,
		"--port", fmt.Sprintf("%d", puerto),
	)
	cmd.Dir = tienda.Ruta

	// Crear el servidor activo antes de iniciar
	servidor := &ServidorActivo{
		Tienda:   tienda,
		Proceso:  cmd,
		Puerto:   puerto,
		Iniciado: time.Now(),
		URL:      fmt.Sprintf("http://127.0.0.1:%d", puerto),
		Activo:   true,
		Logs:     make([]string, 0),
	}

	// Capturar stdin para poder enviar input
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stdin: %v", err)
	}
	servidor.Stdin = stdin

	// Capturar stdout y stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stderr: %v", err)
	}

	// Iniciar el proceso en background
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error al iniciar servidor: %v", err)
	}

	// Registrar el servidor y puerto
	g.servidores[tienda.Nombre] = servidor
	g.puertos[puerto] = true

	// Goroutine para leer stdout
	go func() {
		leerLogs(stdout, servidor)
	}()

	// Goroutine para leer stderr
	go func() {
		leerLogs(stderr, servidor)
	}()

	// Goroutine para detectar cuando el proceso termina
	go func() {
		cmd.Wait()
		g.mutex.Lock()
		defer g.mutex.Unlock()
		if s, ok := g.servidores[tienda.Nombre]; ok {
			s.Activo = false
			s.AgregarLog("--- Servidor detenido ---")
			delete(g.puertos, s.Puerto)
		}
	}()

	return servidor, nil
}

// leerLogs lee de un pipe y agrega las líneas al servidor
func leerLogs(pipe io.Reader, servidor *ServidorActivo) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		servidor.AgregarLog(scanner.Text())
	}
}

// DetenerServidor detiene un servidor específico
func (g *GestorServidores) DetenerServidor(nombreTienda string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	servidor, existe := g.servidores[nombreTienda]
	if !existe {
		return fmt.Errorf("no hay servidor para '%s'", nombreTienda)
	}

	if !servidor.Activo {
		return fmt.Errorf("el servidor de '%s' ya está detenido", nombreTienda)
	}

	// Matar el proceso
	if servidor.Proceso != nil && servidor.Proceso.Process != nil {
		if err := servidor.Proceso.Process.Kill(); err != nil {
			return fmt.Errorf("error al detener servidor: %v", err)
		}
	}

	servidor.Activo = false
	delete(g.puertos, servidor.Puerto)

	return nil
}

// DetenerTodos detiene todos los servidores activos
func (g *GestorServidores) DetenerTodos() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, servidor := range g.servidores {
		if servidor.Activo && servidor.Proceso != nil && servidor.Proceso.Process != nil {
			servidor.Proceso.Process.Kill()
			servidor.Activo = false
		}
	}
	g.puertos = make(map[int]bool)
}

// ObtenerServidoresActivos retorna una lista de servidores que están corriendo
func (g *GestorServidores) ObtenerServidoresActivos() []*ServidorActivo {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	var activos []*ServidorActivo
	for _, servidor := range g.servidores {
		if servidor.Activo {
			activos = append(activos, servidor)
		}
	}
	return activos
}

// ContarActivos retorna la cantidad de servidores activos
func (g *GestorServidores) ContarActivos() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	count := 0
	for _, servidor := range g.servidores {
		if servidor.Activo {
			count++
		}
	}
	return count
}

// TieneServidorActivo verifica si una tienda tiene un servidor corriendo
func (g *GestorServidores) TieneServidorActivo(nombreTienda string) bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe {
		return servidor.Activo
	}
	return false
}

// ObtenerServidor retorna el servidor de una tienda específica
func (g *GestorServidores) ObtenerServidor(nombreTienda string) *ServidorActivo {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe && servidor.Activo {
		return servidor
	}
	return nil
}
