package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/JacuXx/shopify-cli/internal/domain"
)

// ProcessManager implementa Manager usando procesos reales del sistema operativo.
type ProcessManager struct {
	servidores map[string]*ServidorActivo
	mutex      sync.RWMutex
	puertos    map[int]bool
	procesos   map[string]*exec.Cmd
}

// NewProcessManager crea un gestor de servidores listo para usar.
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		servidores: make(map[string]*ServidorActivo),
		puertos:    make(map[int]bool),
		procesos:   make(map[string]*exec.Cmd),
	}
}

// IniciarServidor lanza shopify theme dev para la tienda dada.
func (g *ProcessManager) IniciarServidor(tienda domain.Tienda) (*ServidorActivo, error) {
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
	g.procesos[tienda.Nombre] = cmd
	g.puertos[puerto] = true

	go leerLogs(stdout, servidor)
	go leerLogs(stderr, servidor)

	go func() {
		_ = cmd.Wait() //nolint:errcheck // proceso terminado por señal externa; error no accionable
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

// DetenerServidor detiene el proceso de una tienda específica.
func (g *ProcessManager) DetenerServidor(nombreTienda string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	servidor, existe := g.servidores[nombreTienda]
	if !existe {
		return fmt.Errorf("no hay servidor para '%s'", nombreTienda)
	}

	if !servidor.Activo {
		return fmt.Errorf("el servidor de '%s' ya está detenido", nombreTienda)
	}

	cmd, tieneCmd := g.procesos[nombreTienda]
	if tieneCmd && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("error al detener servidor: %v", err)
		}
	}

	servidor.Activo = false
	delete(g.puertos, servidor.Puerto)

	return nil
}

// DetenerTodos detiene todos los servidores activos.
func (g *ProcessManager) DetenerTodos() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for nombre, servidor := range g.servidores {
		if servidor.Activo {
			if cmd, ok := g.procesos[nombre]; ok && cmd.Process != nil {
				_ = cmd.Process.Kill() //nolint:errcheck // proceso puede ya estar muerto; error esperado
			}
			servidor.Activo = false
		}
	}
	g.puertos = make(map[int]bool)
}

// ObtenerServidoresActivos retorna solo los servidores con Activo=true.
func (g *ProcessManager) ObtenerServidoresActivos() []*ServidorActivo {
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

// ContarActivos retorna el número de servidores activos.
func (g *ProcessManager) ContarActivos() int {
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

// TieneServidorActivo verifica si hay un servidor activo para la tienda.
func (g *ProcessManager) TieneServidorActivo(nombreTienda string) bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe {
		return servidor.Activo
	}
	return false
}

// ObtenerServidor retorna el servidor activo para la tienda, o nil si no existe.
func (g *ProcessManager) ObtenerServidor(nombreTienda string) *ServidorActivo {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if servidor, existe := g.servidores[nombreTienda]; existe && servidor.Activo {
		return servidor
	}
	return nil
}

func leerLogs(pipe io.Reader, servidor *ServidorActivo) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		servidor.AgregarLog(scanner.Text())
	}
}
