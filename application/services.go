package application

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/JacuXx/shopify-cli/domain"
)

// ServerManager gestiona servidores activos
type ServerManager struct {
	servers map[string]*domain.ServidorActivo
	ports   map[int]bool
	mu      sync.RWMutex
}

// NewServerManager crea un nuevo gestor de servidores
func NewServerManager() *ServerManager {
	return &ServerManager{
		servers: make(map[string]*domain.ServidorActivo),
		ports:   make(map[int]bool),
	}
}

// StartServer inicia un servidor de desarrollo
func (sm *ServerManager) StartServer(ctx context.Context, tienda domain.Tienda) (*domain.ServidorActivo, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Verificar si ya existe un servidor activo
	if srv, exists := sm.servers[tienda.Nombre]; exists && srv.Activo {
		return nil, fmt.Errorf("server already running for %s", tienda.Nombre)
	}

	// Obtener puerto disponible
	puerto := sm.getAvailablePort()

	// Crear comando
	cmd := exec.CommandContext(ctx,
		"shopify", "theme", "dev",
		"--store", tienda.URL,
		"--port", fmt.Sprintf("%d", puerto),
	)
	cmd.Dir = tienda.Ruta

	// Crear servidor
	servidor := &domain.ServidorActivo{
		ID:       fmt.Sprintf("%s-%d", tienda.Nombre, time.Now().Unix()),
		Tienda:   tienda,
		Puerto:   puerto,
		Iniciado: time.Now(),
		URL:      fmt.Sprintf("http://127.0.0.1:%d", puerto),
		Activo:   true,
		Logs:     make([]string, 0),
	}

	// Capturar stdout y stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to capture stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to capture stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	sm.servers[tienda.Nombre] = servidor
	sm.ports[puerto] = true

	// Leer logs en background
	go sm.readLogs(stdout, servidor)
	go sm.readLogs(stderr, servidor)

	// Esperar a que el proceso termine
	go func() {
		cmd.Wait()
		sm.mu.Lock()
		defer sm.mu.Unlock()
		if s, ok := sm.servers[tienda.Nombre]; ok {
			s.Activo = false
			s.Logs = append(s.Logs, "--- Server stopped ---")
			delete(sm.ports, s.Puerto)
		}
	}()

	return servidor, nil
}

// StopServer detiene un servidor
func (sm *ServerManager) StopServer(tiendaNombre string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	servidor, exists := sm.servers[tiendaNombre]
	if !exists {
		return fmt.Errorf("no server found for %s", tiendaNombre)
	}

	if !servidor.Activo {
		return fmt.Errorf("server for %s is not active", tiendaNombre)
	}

	servidor.Activo = false
	delete(sm.ports, servidor.Puerto)
	return nil
}

// GetActiveServers retorna todos los servidores activos
func (sm *ServerManager) GetActiveServers() []*domain.ServidorActivo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var activos []*domain.ServidorActivo
	for _, srv := range sm.servers {
		if srv.Activo {
			activos = append(activos, srv)
		}
	}
	return activos
}

// GetServer obtiene un servidor específico
func (sm *ServerManager) GetServer(tiendaNombre string) *domain.ServidorActivo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if srv, ok := sm.servers[tiendaNombre]; ok && srv.Activo {
		return srv
	}
	return nil
}

// HasActiveServer verifica si hay un servidor activo para una tienda
func (sm *ServerManager) HasActiveServer(tiendaNombre string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if srv, ok := sm.servers[tiendaNombre]; ok {
		return srv.Activo
	}
	return false
}

// StopAll detiene todos los servidores
func (sm *ServerManager) StopAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, srv := range sm.servers {
		if srv.Activo {
			srv.Activo = false
		}
	}
	sm.ports = make(map[int]bool)
}

// Private helper methods

func (sm *ServerManager) getAvailablePort() int {
	puerto := 9292
	for sm.ports[puerto] {
		puerto++
	}
	return puerto
}

func (sm *ServerManager) readLogs(reader io.Reader, servidor *domain.ServidorActivo) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			sm.mu.Lock()
			servidor.Logs = append(servidor.Logs, string(buf[:n]))
			if len(servidor.Logs) > 1000 {
				servidor.Logs = servidor.Logs[len(servidor.Logs)-1000:]
			}
			sm.mu.Unlock()
		}
	}
}
