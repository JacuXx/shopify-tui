package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type ServidorActivo struct {
	Tienda    Tienda
	Proceso   *exec.Cmd
	Puerto    int
	Iniciado  time.Time
	URL       string
	Activo    bool
	Logs      []string
	LogsMutex sync.RWMutex
	Stdin     io.WriteCloser
}

func (s *ServidorActivo) AgregarLog(linea string) {
	s.LogsMutex.Lock()
	defer s.LogsMutex.Unlock()

	s.Logs = append(s.Logs, linea)

	if len(s.Logs) > 100 {
		s.Logs = s.Logs[len(s.Logs)-100:]
	}
}

func (s *ServidorActivo) ObtenerLogs() []string {
	s.LogsMutex.RLock()
	defer s.LogsMutex.RUnlock()

	copia := make([]string, len(s.Logs))
	copy(copia, s.Logs)
	return copia
}

func (s *ServidorActivo) EnviarInput(input string) error {
	if s.Stdin == nil {
		return fmt.Errorf("stdin no disponible")
	}
	_, err := s.Stdin.Write([]byte(input))
	return err
}

type GestorServidores struct {
	servidores map[string]*ServidorActivo
	mutex      sync.RWMutex
	puertos    map[int]bool
}

var gestorGlobal = &GestorServidores{
	servidores: make(map[string]*ServidorActivo),
	puertos:    make(map[int]bool),
}

func ObtenerGestor() *GestorServidores {
	return gestorGlobal
}

func (g *GestorServidores) ObtenerPuertoDisponible() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	puerto := 9292
	for g.puertos[puerto] {
		puerto++
	}
	return puerto
}

func (g *GestorServidores) IniciarServidor(tienda Tienda) (*ServidorActivo, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if servidor, existe := g.servidores[tienda.Nombre]; existe && servidor.Activo {
		return nil, fmt.Errorf("ya hay un servidor activo para '%s'", tienda.Nombre)
	}

	puerto := 9292
	for g.puertos[puerto] {
		puerto++
	}

	cmd := exec.Command("shopify", "theme", "dev",
		"--store", tienda.URL,
		"--port", fmt.Sprintf("%d", puerto),
	)
	cmd.Dir = tienda.Ruta

	servidor := &ServidorActivo{
		Tienda:   tienda,
		Proceso:  cmd,
		Puerto:   puerto,
		Iniciado: time.Now(),
		URL:      fmt.Sprintf("http://127.0.0.1:%d", puerto),
		Activo:   true,
		Logs:     make([]string, 0),
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stdin: %v", err)
	}
	servidor.Stdin = stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("error al capturar stderr: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error al iniciar servidor: %v", err)
	}

	g.servidores[tienda.Nombre] = servidor
	g.puertos[puerto] = true

	go func() {
		leerLogs(stdout, servidor)
	}()

	go func() {
		leerLogs(stderr, servidor)
	}()

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

func leerLogs(pipe io.Reader, servidor *ServidorActivo) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		servidor.AgregarLog(scanner.Text())
	}
}

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

	if servidor.Proceso != nil && servidor.Proceso.Process != nil {
		if err := servidor.Proceso.Process.Kill(); err != nil {
			return fmt.Errorf("error al detener servidor: %v", err)
		}
	}

	servidor.Activo = false
	delete(g.puertos, servidor.Puerto)

	return nil
}

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

func (g *GestorServidores) TieneServidorActivo(nombreTienda string) bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe {
		return servidor.Activo
	}
	return false
}

func (g *GestorServidores) ObtenerServidor(nombreTienda string) *ServidorActivo {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe && servidor.Activo {
		return servidor
	}
	return nil
}
